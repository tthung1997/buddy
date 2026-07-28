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
// than failing the whole match. Staples are satisfied by presence: a recipe
// still says how much salt to add, but the pantry only records whether there
// is any.
type PantryCookEngine struct {
	Converter cooking.IUnitConverter
}

func NewPantryCookEngine(converter cooking.IUnitConverter) *PantryCookEngine {
	if converter == nil {
		converter = NewSimpleUnitConverter()
	}
	return &PantryCookEngine{Converter: converter}
}

func (e *PantryCookEngine) Match(recipes []cooking.Recipe, pantry cooking.Pantry, options cooking.MatchOptions) []cooking.RecipeMatch {
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

func (e *PantryCookEngine) matchRecipe(recipe cooking.Recipe, pantry cooking.Pantry, options cooking.MatchOptions) cooking.RecipeMatch {
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
	for _, ingredient := range e.aggregate(target.RequiredIngredients(), pantry) {
		blocking := !ingredient.Optional || options.IncludeOptional
		if blocking {
			required++
		}

		enough := false
		if pantry.IsStaple(ingredient.IngredientId) {
			enough = e.matchStaple(&match, ingredient, pantry)
		} else {
			enough = e.matchMeasured(&match, ingredient, pantry)
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

// matchStaple satisfies an ingredient on presence alone. The recipe amount is
// still carried through so the page can show how much to add.
func (e *PantryCookEngine) matchStaple(match *cooking.RecipeMatch, ingredient cooking.RecipeIngredient, pantry cooking.Pantry) bool {
	stock := pantry.StockOf(ingredient.IngredientId)
	if !stock.Available() {
		match.Missing = append(match.Missing, cooking.MissingIngredient{
			IngredientId: ingredient.IngredientId,
			Optional:     ingredient.Optional,
			Staple:       true,
			Reason:       cooking.MissingReasonOutOfStock,
		})
		return false
	}

	match.Have = append(match.Have, cooking.MatchedIngredient{
		IngredientId: ingredient.IngredientId,
		Required:     ingredient.Amount,
		Unit:         ingredient.Unit,
		Optional:     ingredient.Optional,
		Staple:       true,
		Stock:        stock,
	})
	return true
}

func (e *PantryCookEngine) matchMeasured(match *cooking.RecipeMatch, ingredient cooking.RecipeIngredient, pantry cooking.Pantry) bool {
	available, present, compatible := e.availableAmount(pantry.Items, ingredient.IngredientId, ingredient.Unit)

	switch {
	case ingredient.Amount <= 0:
		match.Have = append(match.Have, matchedIngredient(ingredient, available))
		return true
	case !present:
		match.Missing = append(match.Missing, missingIngredient(ingredient, ingredient.Amount, cooking.MissingReasonAbsent))
		return false
	case !compatible:
		match.Missing = append(match.Missing, missingIngredient(ingredient, ingredient.Amount, cooking.MissingReasonUnitMismatch))
		return false
	case available+amountEpsilon >= ingredient.Amount:
		match.Have = append(match.Have, matchedIngredient(ingredient, available))
		return true
	default:
		match.Have = append(match.Have, matchedIngredient(ingredient, available))
		match.Missing = append(match.Missing, missingIngredient(ingredient, ingredient.Amount-available, cooking.MissingReasonInsufficient))
		return false
	}
}

func (e *PantryCookEngine) ShoppingList(plans []cooking.RecipePlan, pantry cooking.Pantry) []cooking.MissingIngredient {
	remaining := clonePantry(pantry.Items)
	missing := []cooking.MissingIngredient{}

	for _, plan := range plans {
		target := plan.Recipe
		if plan.Servings > 0 {
			if scaled, err := plan.Recipe.ScaleTo(plan.Servings); err == nil {
				target = scaled
			}
		}

		for _, ingredient := range e.aggregate(target.RequiredIngredients(), pantry) {
			if ingredient.Optional {
				continue
			}

			// A staple goes on the list when it has run out or is running low,
			// and never carries a quantity because none was ever tracked.
			if pantry.IsStaple(ingredient.IngredientId) {
				stock := pantry.StockOf(ingredient.IngredientId)
				if stock == cooking.InStock {
					continue
				}
				reason := cooking.MissingReasonOutOfStock
				if stock == cooking.LowStock {
					reason = cooking.MissingReasonRunningLow
				}
				missing = e.addStaple(missing, ingredient.IngredientId, reason)
				continue
			}

			if ingredient.Amount <= 0 {
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

func (e *PantryCookEngine) Consume(pantry cooking.Pantry, recipe cooking.Recipe, servings int) ([]cooking.PantryItem, error) {
	target := recipe
	if servings > 0 && servings != recipe.Servings {
		scaled, err := recipe.ScaleTo(servings)
		if err != nil {
			return nil, err
		}
		target = scaled
	}

	updated := clonePantry(pantry.Items)
	now := time.Now()

	for _, ingredient := range e.aggregate(target.RequiredIngredients(), pantry) {
		if ingredient.Optional {
			continue
		}

		// Staples are never deducted; there is no amount to deduct from. They
		// still have to be in stock for the recipe to be cooked at all.
		if pantry.IsStaple(ingredient.IngredientId) {
			if !pantry.StockOf(ingredient.IngredientId).Available() {
				return nil, fmt.Errorf("%w: %s has run out", cooking.ErrInsufficientPantry, ingredient.IngredientId)
			}
			continue
		}

		if ingredient.Amount <= 0 {
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
		// Staple rows carry a level rather than an amount, so they survive even
		// though their amount is zero.
		if item.Stock.IsValid() || pantry.IsStaple(item.IngredientId) || item.Amount > amountEpsilon {
			pruned = append(pruned, item)
		}
	}

	return pruned, nil
}

// aggregate merges repeated mentions of the same ingredient into a single
// requirement, converting into the unit of the first mention when possible.
// Staples are merged by identity alone since their amounts are not tracked.
func (e *PantryCookEngine) aggregate(ingredients []cooking.RecipeIngredient, pantry cooking.Pantry) []cooking.RecipeIngredient {
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
			if pantry.IsStaple(ingredient.IngredientId) {
				aggregated[i].Optional = aggregated[i].Optional && ingredient.Optional
				merged = true
				break
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
func (e *PantryCookEngine) availableAmount(items []cooking.PantryItem, ingredientId string, unit cooking.Unit) (float64, bool, bool) {
	total := 0.0
	present := false
	compatible := false

	for _, item := range items {
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
func (e *PantryCookEngine) drawDown(items []cooking.PantryItem, ingredient cooking.RecipeIngredient) (float64, bool, bool) {
	needed := ingredient.Amount
	present := false
	compatible := false

	for i := range items {
		if items[i].IngredientId != ingredient.IngredientId {
			continue
		}
		present = true

		availableInUnit, err := e.Converter.Convert(items[i].Amount, items[i].Unit, ingredient.Unit)
		if err != nil {
			continue
		}
		compatible = true

		used := math.Min(availableInUnit, needed)
		needed -= used
		if leftover, err := e.Converter.Convert(availableInUnit-used, ingredient.Unit, items[i].Unit); err == nil {
			items[i].Amount = leftover
		}
		if needed <= amountEpsilon {
			break
		}
	}

	return needed, present, compatible
}

func (e *PantryCookEngine) deduct(items []cooking.PantryItem, ingredient cooking.RecipeIngredient, at time.Time) {
	needed := ingredient.Amount
	for i := range items {
		if items[i].IngredientId != ingredient.IngredientId || needed <= amountEpsilon {
			continue
		}

		availableInUnit, err := e.Converter.Convert(items[i].Amount, items[i].Unit, ingredient.Unit)
		if err != nil {
			continue
		}

		used := math.Min(availableInUnit, needed)
		needed -= used
		if leftover, err := e.Converter.Convert(availableInUnit-used, ingredient.Unit, items[i].Unit); err == nil {
			items[i].Amount = leftover
			items[i].UpdatedDateTime = at
		}
	}
}

func (e *PantryCookEngine) addMissing(missing []cooking.MissingIngredient, item cooking.MissingIngredient) []cooking.MissingIngredient {
	for i := range missing {
		if missing[i].IngredientId != item.IngredientId || missing[i].Staple {
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

// addStaple records a staple once, preferring the more urgent reason when the
// same staple is needed by several planned recipes.
func (e *PantryCookEngine) addStaple(missing []cooking.MissingIngredient, ingredientId string, reason cooking.MissingReason) []cooking.MissingIngredient {
	for i := range missing {
		if missing[i].IngredientId != ingredientId || !missing[i].Staple {
			continue
		}
		if reason == cooking.MissingReasonOutOfStock {
			missing[i].Reason = reason
		}
		return missing
	}
	return append(missing, cooking.MissingIngredient{
		IngredientId: ingredientId,
		Staple:       true,
		Reason:       reason,
	})
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

func clonePantry(items []cooking.PantryItem) []cooking.PantryItem {
	clone := make([]cooking.PantryItem, len(items))
	copy(clone, items)
	return clone
}
