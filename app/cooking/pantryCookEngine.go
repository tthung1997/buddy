package cooking

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/tthung1997/buddy/core/cooking"
)

// amountEpsilon absorbs the rounding noise of unit conversion so a pantry that
// holds exactly what a recipe needs still counts as sufficient.
const amountEpsilon = 1e-9

// PantryCookEngine answers "what can I cook with what I have?" by joining a
// pantry against recipe ingredients. Ingredients stored in a unit that cannot
// be compared with the one a recipe asks for are reported as missing rather
// than failing the whole match.
type PantryCookEngine struct {
	Converter cooking.IUnitConverter
}

func NewPantryCookEngine(converter cooking.IUnitConverter) *PantryCookEngine {
	if converter == nil {
		converter = NewSimpleUnitConverter()
	}
	return &PantryCookEngine{Converter: converter}
}

func (e *PantryCookEngine) Match(recipes []cooking.Recipe, pantry []cooking.PantryItem, options cooking.MatchOptions) []cooking.RecipeMatch {
	matches := []cooking.RecipeMatch{}
	for _, recipe := range recipes {
		match := e.matchRecipe(recipe, pantry, options)
		if options.OnlyCookable && !match.Cookable {
			continue
		}
		if match.Coverage+amountEpsilon < options.MinCoverage {
			continue
		}
		matches = append(matches, match)
	}

	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].Cookable != matches[j].Cookable {
			return matches[i].Cookable
		}
		if math.Abs(matches[i].Coverage-matches[j].Coverage) > amountEpsilon {
			return matches[i].Coverage > matches[j].Coverage
		}
		return matches[i].Recipe.Name < matches[j].Recipe.Name
	})

	return matches
}

func (e *PantryCookEngine) matchRecipe(recipe cooking.Recipe, pantry []cooking.PantryItem, options cooking.MatchOptions) cooking.RecipeMatch {
	target := recipe
	servings := recipe.Servings
	if options.Servings > 0 {
		if scaled, err := recipe.ScaleTo(options.Servings); err == nil {
			target = scaled
			servings = options.Servings
		}
	}

	match := cooking.RecipeMatch{
		Recipe:   target,
		Servings: servings,
		Have:     []cooking.MatchedIngredient{},
		Missing:  []cooking.MissingIngredient{},
		Coverage: 1,
		Cookable: true,
	}

	required := 0
	satisfied := 0
	for _, ingredient := range e.aggregate(target.RequiredIngredients()) {
		blocking := !ingredient.Optional || options.IncludeOptional
		if blocking {
			required++
		}

		available, present, compatible := e.availableAmount(pantry, ingredient.IngredientId, ingredient.Unit)
		enough := false
		switch {
		case ingredient.Amount <= 0:
			enough = true
			match.Have = append(match.Have, matchedIngredient(ingredient, available))
		case !present:
			match.Missing = append(match.Missing, missingIngredient(ingredient, ingredient.Amount, cooking.MissingReasonAbsent))
		case !compatible:
			match.Missing = append(match.Missing, missingIngredient(ingredient, ingredient.Amount, cooking.MissingReasonUnitMismatch))
		case available+amountEpsilon >= ingredient.Amount:
			enough = true
			match.Have = append(match.Have, matchedIngredient(ingredient, available))
		default:
			match.Have = append(match.Have, matchedIngredient(ingredient, available))
			match.Missing = append(match.Missing, missingIngredient(ingredient, ingredient.Amount-available, cooking.MissingReasonInsufficient))
		}

		if !blocking {
			continue
		}
		if enough {
			satisfied++
		} else {
			match.Cookable = false
		}
	}

	if required > 0 {
		match.Coverage = float64(satisfied) / float64(required)
	}

	return match
}

func (e *PantryCookEngine) ShoppingList(plans []cooking.RecipePlan, pantry []cooking.PantryItem) []cooking.MissingIngredient {
	remaining := clonePantry(pantry)
	missing := []cooking.MissingIngredient{}

	for _, plan := range plans {
		target := plan.Recipe
		if plan.Servings > 0 {
			if scaled, err := plan.Recipe.ScaleTo(plan.Servings); err == nil {
				target = scaled
			}
		}

		for _, ingredient := range e.aggregate(target.RequiredIngredients()) {
			if ingredient.Optional || ingredient.Amount <= 0 {
				continue
			}

			needed, present, compatible := e.drawDown(remaining, ingredient)
			if needed <= amountEpsilon {
				continue
			}

			reason := cooking.MissingReasonAbsent
			if present && compatible {
				reason = cooking.MissingReasonInsufficient
			} else if present {
				reason = cooking.MissingReasonUnitMismatch
			}
			missing = e.addMissing(missing, missingIngredient(ingredient, needed, reason))
		}
	}

	sort.SliceStable(missing, func(i, j int) bool {
		if missing[i].IngredientId != missing[j].IngredientId {
			return missing[i].IngredientId < missing[j].IngredientId
		}
		return missing[i].Unit < missing[j].Unit
	})

	return missing
}

func (e *PantryCookEngine) Consume(pantry []cooking.PantryItem, recipe cooking.Recipe, servings int) ([]cooking.PantryItem, error) {
	target := recipe
	if servings > 0 && servings != recipe.Servings {
		scaled, err := recipe.ScaleTo(servings)
		if err != nil {
			return nil, err
		}
		target = scaled
	}

	updated := clonePantry(pantry)
	now := time.Now()

	for _, ingredient := range e.aggregate(target.RequiredIngredients()) {
		if ingredient.Optional || ingredient.Amount <= 0 {
			continue
		}

		available, present, compatible := e.availableAmount(updated, ingredient.IngredientId, ingredient.Unit)
		if !present {
			return nil, fmt.Errorf("%w: %s is not in the pantry", cooking.ErrInsufficientPantry, ingredient.IngredientId)
		}
		if !compatible {
			return nil, fmt.Errorf("%w: %s is stored in a unit that cannot be measured in %s",
				cooking.ErrIncompatibleUnits, ingredient.IngredientId, ingredient.Unit)
		}
		if available+amountEpsilon < ingredient.Amount {
			return nil, fmt.Errorf("%w: %s needs %g %s but only %g %s is available",
				cooking.ErrInsufficientPantry, ingredient.IngredientId, ingredient.Amount, ingredient.Unit, available, ingredient.Unit)
		}

		e.deduct(updated, ingredient, now)
	}

	pruned := []cooking.PantryItem{}
	for _, item := range updated {
		if item.Amount > amountEpsilon {
			pruned = append(pruned, item)
		}
	}

	return pruned, nil
}

// aggregate merges repeated mentions of the same ingredient into a single
// requirement, converting into the unit of the first mention when possible.
func (e *PantryCookEngine) aggregate(ingredients []cooking.RecipeIngredient) []cooking.RecipeIngredient {
	aggregated := []cooking.RecipeIngredient{}
	for _, ingredient := range ingredients {
		if ingredient.IngredientId == "" {
			continue
		}

		merged := false
		for i := range aggregated {
			if aggregated[i].IngredientId != ingredient.IngredientId {
				continue
			}
			amount, err := e.Converter.Convert(ingredient.Amount, ingredient.Unit, aggregated[i].Unit)
			if err != nil {
				continue
			}
			aggregated[i].Amount += amount
			aggregated[i].Optional = aggregated[i].Optional && ingredient.Optional
			merged = true
			break
		}

		if !merged {
			aggregated = append(aggregated, ingredient)
		}
	}
	return aggregated
}

// availableAmount totals how much of an ingredient the pantry holds, expressed
// in the requested unit. It also reports whether the ingredient is in the
// pantry at all and whether any entry could be converted.
func (e *PantryCookEngine) availableAmount(pantry []cooking.PantryItem, ingredientId string, unit cooking.Unit) (float64, bool, bool) {
	total := 0.0
	present := false
	compatible := false

	for _, item := range pantry {
		if item.IngredientId != ingredientId {
			continue
		}
		present = true
		amount, err := e.Converter.Convert(item.Amount, item.Unit, unit)
		if err != nil {
			continue
		}
		compatible = true
		total += amount
	}

	return total, present, compatible
}

// drawDown removes as much of the requirement as the pantry can cover and
// returns what is still needed.
func (e *PantryCookEngine) drawDown(pantry []cooking.PantryItem, ingredient cooking.RecipeIngredient) (float64, bool, bool) {
	needed := ingredient.Amount
	present := false
	compatible := false

	for i := range pantry {
		if pantry[i].IngredientId != ingredient.IngredientId {
			continue
		}
		present = true

		availableInUnit, err := e.Converter.Convert(pantry[i].Amount, pantry[i].Unit, ingredient.Unit)
		if err != nil {
			continue
		}
		compatible = true

		used := math.Min(availableInUnit, needed)
		needed -= used
		if leftover, err := e.Converter.Convert(availableInUnit-used, ingredient.Unit, pantry[i].Unit); err == nil {
			pantry[i].Amount = leftover
		}
		if needed <= amountEpsilon {
			break
		}
	}

	return needed, present, compatible
}

func (e *PantryCookEngine) deduct(pantry []cooking.PantryItem, ingredient cooking.RecipeIngredient, at time.Time) {
	needed := ingredient.Amount
	for i := range pantry {
		if pantry[i].IngredientId != ingredient.IngredientId || needed <= amountEpsilon {
			continue
		}

		availableInUnit, err := e.Converter.Convert(pantry[i].Amount, pantry[i].Unit, ingredient.Unit)
		if err != nil {
			continue
		}

		used := math.Min(availableInUnit, needed)
		needed -= used
		if leftover, err := e.Converter.Convert(availableInUnit-used, ingredient.Unit, pantry[i].Unit); err == nil {
			pantry[i].Amount = leftover
			pantry[i].UpdatedDateTime = at
		}
	}
}

func (e *PantryCookEngine) addMissing(missing []cooking.MissingIngredient, item cooking.MissingIngredient) []cooking.MissingIngredient {
	for i := range missing {
		if missing[i].IngredientId != item.IngredientId {
			continue
		}
		amount, err := e.Converter.Convert(item.Amount, item.Unit, missing[i].Unit)
		if err != nil {
			continue
		}
		missing[i].Amount += amount
		return missing
	}
	return append(missing, item)
}

func matchedIngredient(ingredient cooking.RecipeIngredient, available float64) cooking.MatchedIngredient {
	return cooking.MatchedIngredient{
		IngredientId: ingredient.IngredientId,
		Required:     ingredient.Amount,
		Available:    available,
		Unit:         ingredient.Unit,
		Optional:     ingredient.Optional,
	}
}

func missingIngredient(ingredient cooking.RecipeIngredient, amount float64, reason cooking.MissingReason) cooking.MissingIngredient {
	return cooking.MissingIngredient{
		IngredientId: ingredient.IngredientId,
		Amount:       amount,
		Unit:         ingredient.Unit,
		Optional:     ingredient.Optional,
		Reason:       reason,
	}
}

func clonePantry(pantry []cooking.PantryItem) []cooking.PantryItem {
	clone := make([]cooking.PantryItem, len(pantry))
	copy(clone, pantry)
	return clone
}
