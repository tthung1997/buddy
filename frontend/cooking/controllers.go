package cooking

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/tthung1997/buddy/core/cooking"
	"github.com/tthung1997/buddy/core/random"
	"github.com/tthung1997/buddy/frontend/shopping"
)

const cookingStaticDir = "frontend/cooking/static"

var (
	cookingIndexTmpl   = template.Must(template.ParseFiles(cookingStaticDir + "/index.html"))
	cookingRecipesTmpl = template.Must(template.ParseFiles(cookingStaticDir + "/recipes.html"))
	cookingRecipeTmpl  = template.Must(template.ParseFiles(cookingStaticDir + "/recipe.html"))
	cookingCookTmpl    = template.Must(template.ParseFiles(cookingStaticDir + "/cook.html"))
	cookingPantryTmpl  = template.Must(template.ParseFiles(cookingStaticDir + "/pantry.html"))
	cookingEditTmpl    = template.Must(template.ParseFiles(cookingStaticDir + "/edit.html"))
)

type CookingController struct {
	RecipeRepository     cooking.IRecipeRepository
	IngredientRepository cooking.IIngredientRepository
	PantryRepository     cooking.IPantryRepository
	Engine               cooking.ICookEngine
	Randomizer           random.IRandomizer
}

func NewCookingController(
	recipes cooking.IRecipeRepository,
	ingredients cooking.IIngredientRepository,
	pantry cooking.IPantryRepository,
	engine cooking.ICookEngine,
	randomizer random.IRandomizer,
) *CookingController {
	return &CookingController{
		RecipeRepository:     recipes,
		IngredientRepository: ingredients,
		PantryRepository:     pantry,
		Engine:               engine,
		Randomizer:           randomizer,
	}
}

// Index renders the cooking dashboard.
func (c *CookingController) Index(w http.ResponseWriter, r *http.Request) {
	data := IndexPageData{
		Success:    r.URL.Query().Get("success"),
		Error:      r.URL.Query().Get("error"),
		Highlights: []MatchView{},
		NextUp:     []MissingLine{},
	}

	recipes, index, pantry, err := c.load()
	if err != nil {
		data.Error = err.Error()
		render(w, cookingIndexTmpl, data)
		return
	}

	data.RecipeCount = len(recipes)
	data.IngredientCount = len(index)
	data.PantryCount = len(pantry)

	matches := c.Engine.Match(recipes, pantry, cooking.MatchOptions{})
	for _, match := range matches {
		if match.Cookable {
			data.CookableCount++
		}
	}

	for _, match := range matches {
		if len(data.Highlights) == 3 {
			break
		}
		data.Highlights = append(data.Highlights, buildMatchView(match, index))
	}

	for _, match := range matches {
		if match.Cookable {
			continue
		}
		view := buildMatchView(match, index)
		for _, missing := range view.Missing {
			if missing.Optional || len(data.NextUp) == 5 {
				continue
			}
			data.NextUp = append(data.NextUp, missing)
		}
		break
	}

	render(w, cookingIndexTmpl, data)
}

// Recipes renders the recipe library with optional search and tag filters.
func (c *CookingController) Recipes(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	data := RecipesPageData{
		Success: params.Get("success"),
		Error:   params.Get("error"),
		Query:   strings.TrimSpace(params.Get("q")),
		Tag:     strings.TrimSpace(params.Get("tag")),
		Recipes: []RecipeView{},
		Tags:    []string{},
	}

	recipes, err := c.RecipeRepository.List()
	if err != nil {
		data.Error = err.Error()
		render(w, cookingRecipesTmpl, data)
		return
	}
	index, err := c.ingredientIndex()
	if err != nil {
		data.Error = err.Error()
		render(w, cookingRecipesTmpl, data)
		return
	}

	data.Total = len(recipes)
	tags := map[string]bool{}
	query := strings.ToLower(data.Query)

	for _, recipe := range recipes {
		for _, tag := range recipe.Tags {
			tags[tag] = true
		}
		if !matchesQuery(recipe, query) || !hasTag(recipe, data.Tag) {
			continue
		}
		data.Recipes = append(data.Recipes, buildRecipeView(recipe, index))
	}

	for tag := range tags {
		data.Tags = append(data.Tags, tag)
	}
	sort.Strings(data.Tags)
	data.Count = len(data.Recipes)

	render(w, cookingRecipesTmpl, data)
}

// Recipe renders a single recipe, scaled to the requested serving count and
// matched against the pantry.
func (c *CookingController) Recipe(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	recipe, err := c.RecipeRepository.Get(params.Get("id"))
	if err != nil {
		redirectWithError(w, r, "/cooking/recipes", err)
		return
	}

	index, err := c.ingredientIndex()
	if err != nil {
		redirectWithError(w, r, "/cooking/recipes", err)
		return
	}
	pantry, err := c.PantryRepository.List()
	if err != nil {
		redirectWithError(w, r, "/cooking/recipes", err)
		return
	}

	servings := parseInt(params.Get("servings"), recipe.Servings)
	if servings <= 0 {
		servings = 1
	}

	data := RecipePageData{
		Success:  params.Get("success"),
		Error:    params.Get("error"),
		Servings: servings,
		Recipe:   buildRecipeView(recipe, index),
	}

	matches := c.Engine.Match([]cooking.Recipe{recipe}, pantry, cooking.MatchOptions{Servings: servings})
	if len(matches) > 0 {
		data.Match = buildMatchView(matches[0], index)
		data.Recipe = data.Match.Recipe
		data.HasMatch = true
	}

	render(w, cookingRecipeTmpl, data)
}

// Cook renders the "what can I cook right now" view.
func (c *CookingController) Cook(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	data := CookPageData{
		Success:         params.Get("success"),
		Error:           params.Get("error"),
		Servings:        parseInt(params.Get("servings"), 0),
		IncludeOptional: params.Get("includeOptional") == "on",
		OnlyCookable:    params.Get("onlyCookable") == "on",
		Matches:         []MatchView{},
	}

	recipes, index, pantry, err := c.load()
	if err != nil {
		data.Error = err.Error()
		render(w, cookingCookTmpl, data)
		return
	}

	all := c.Engine.Match(recipes, pantry, cooking.MatchOptions{
		Servings:        data.Servings,
		IncludeOptional: data.IncludeOptional,
	})
	data.TotalCount = len(all)

	for _, match := range all {
		if match.Cookable {
			data.CookableCount++
		}
		if data.OnlyCookable && !match.Cookable {
			continue
		}
		data.Matches = append(data.Matches, buildMatchView(match, index))
	}

	render(w, cookingCookTmpl, data)
}

// Pick picks a random cookable recipe, favouring the ones cooked least recently.
func (c *CookingController) Pick(w http.ResponseWriter, r *http.Request) {
	recipes, _, pantry, err := c.load()
	if err != nil {
		redirectWithError(w, r, "/cooking/cook", err)
		return
	}

	matches := c.Engine.Match(recipes, pantry, cooking.MatchOptions{OnlyCookable: true})
	if len(matches) == 0 {
		http.Redirect(w, r, "/cooking/cook?error="+url.QueryEscape("Nothing can be cooked with the current pantry"), http.StatusSeeOther)
		return
	}

	choices := make([]random.Choice, 0, len(matches))
	for _, match := range matches {
		choices = append(choices, random.Choice{
			Id:     match.Recipe.Id,
			Value:  match.Recipe.Id,
			Weight: stalenessWeight(match.Recipe),
		})
	}

	picked := c.Randomizer.GetChoice(choices, 1)
	if len(picked) == 0 {
		redirectWithError(w, r, "/cooking/cook", fmt.Errorf("could not pick a recipe"))
		return
	}

	http.Redirect(w, r, "/cooking/recipes/view?id="+url.QueryEscape(picked[0].Value)+
		"&success="+url.QueryEscape("Buddy picked this one for you"), http.StatusSeeOther)
}

// Complete deducts a cooked recipe from the pantry.
func (c *CookingController) Complete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/cooking/cook", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	back := "/cooking/recipes/view?id=" + url.QueryEscape(id)

	recipe, err := c.RecipeRepository.Get(id)
	if err != nil {
		redirectWithError(w, r, "/cooking/cook", err)
		return
	}

	servings := parseInt(r.FormValue("servings"), recipe.Servings)
	pantry, err := c.PantryRepository.List()
	if err != nil {
		redirectWithError(w, r, back, err)
		return
	}

	updated, err := c.Engine.Consume(pantry, recipe, servings)
	if err != nil {
		redirectWithError(w, r, back, err)
		return
	}
	if err := c.savePantry(pantry, updated); err != nil {
		redirectWithError(w, r, back, err)
		return
	}

	now := time.Now()
	recipe.LastCookedDateTime = &now
	recipe.UpdatedDateTime = now
	if err := c.RecipeRepository.CreateOrUpdate(recipe); err != nil {
		redirectWithError(w, r, back, err)
		return
	}

	http.Redirect(w, r, back+"&success="+url.QueryEscape("Cooked. The pantry has been updated."), http.StatusSeeOther)
}

// PushToShopping sends everything a recipe is missing to the shopping list.
func (c *CookingController) PushToShopping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/cooking/cook", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	back := "/cooking/recipes/view?id=" + url.QueryEscape(id)

	recipe, err := c.RecipeRepository.Get(id)
	if err != nil {
		redirectWithError(w, r, "/cooking/cook", err)
		return
	}

	index, err := c.ingredientIndex()
	if err != nil {
		redirectWithError(w, r, back, err)
		return
	}
	pantry, err := c.PantryRepository.List()
	if err != nil {
		redirectWithError(w, r, back, err)
		return
	}

	servings := parseInt(r.FormValue("servings"), recipe.Servings)
	missing := c.Engine.ShoppingList([]cooking.RecipePlan{{Recipe: recipe, Servings: servings}}, pantry)
	if len(missing) == 0 {
		http.Redirect(w, r, back+"&success="+url.QueryEscape("Nothing to buy, the pantry already covers this recipe"), http.StatusSeeOther)
		return
	}

	store := strings.TrimSpace(r.FormValue("store"))
	items := make([]shopping.ShoppingItem, 0, len(missing))
	for _, item := range missing {
		items = append(items, shopping.ShoppingItem{
			Name:  fmt.Sprintf("%s (%s %s)", ingredientName(index, item.IngredientId), formatAmount(item.Amount), item.Unit),
			Store: store,
		})
	}

	added, err := shopping.AppendItems(items)
	if err != nil {
		redirectWithError(w, r, back, err)
		return
	}

	message := fmt.Sprintf("Added %d item(s) to the shopping list", added)
	if added == 0 {
		message = "Those items are already on the shopping list"
	}
	http.Redirect(w, r, back+"&success="+url.QueryEscape(message), http.StatusSeeOther)
}

// Pantry renders the pantry and ingredient catalog workspace.
func (c *CookingController) PantryPage(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	data := PantryPageData{
		Success:       params.Get("success"),
		Error:         params.Get("error"),
		Units:         cooking.Units(),
		IngredientIds: []string{},
	}

	ingredients, err := c.IngredientRepository.List()
	if err != nil {
		data.Error = err.Error()
		render(w, cookingPantryTmpl, data)
		return
	}
	for _, ingredient := range ingredients {
		data.IngredientIds = append(data.IngredientIds, ingredient.Name)
	}

	render(w, cookingPantryTmpl, data)
}

// Edit renders the recipe editor for a new or an existing recipe.
func (c *CookingController) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	recipe := cooking.Recipe{Servings: 2, Ingredients: []cooking.RecipeIngredient{}, Steps: []cooking.Step{}}
	isNew := true

	if id != "" {
		stored, err := c.RecipeRepository.Get(id)
		if err != nil {
			redirectWithError(w, r, "/cooking/recipes", err)
			return
		}
		recipe = stored
		isNew = false
	}

	ingredients, err := c.IngredientRepository.List()
	if err != nil {
		redirectWithError(w, r, "/cooking/recipes", err)
		return
	}

	data := EditPageData{
		Error:         r.URL.Query().Get("error"),
		IsNew:         isNew,
		RecipeId:      recipe.Id,
		RecipeName:    recipe.Name,
		RecipeJSON:    mustJSON(recipeForEditor(recipe, ingredients)),
		CatalogJSON:   mustJSON(ingredients),
		UnitsJSON:     mustJSON(cooking.Units()),
		DurationsJSON: mustJSON(cooking.DurationUnits()),
	}

	render(w, cookingEditTmpl, data)
}

// SaveRecipe stores a recipe posted by the editor as JSON.
func (c *CookingController) SaveRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "only POST is supported")
		return
	}

	var recipe cooking.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		writeJSONError(w, http.StatusBadRequest, "could not read the recipe: "+err.Error())
		return
	}

	if err := c.normalizeRecipe(&recipe); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.RecipeRepository.CreateOrUpdate(recipe); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, map[string]string{"id": recipe.Id})
}

// DeleteRecipe removes a recipe.
func (c *CookingController) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/cooking/recipes", http.StatusSeeOther)
		return
	}

	if err := c.RecipeRepository.Delete(r.FormValue("id")); err != nil {
		redirectWithError(w, r, "/cooking/recipes", err)
		return
	}

	http.Redirect(w, r, "/cooking/recipes?success="+url.QueryEscape("Recipe deleted"), http.StatusSeeOther)
}

// IngredientsAPI reads and replaces the ingredient catalog.
func (c *CookingController) IngredientsAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ingredients, err := c.IngredientRepository.List()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, ingredients)
	case http.MethodPost:
		var incoming []cooking.Ingredient
		if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		now := time.Now()
		catalog := []cooking.Ingredient{}
		taken := map[string]bool{}
		for _, ingredient := range incoming {
			ingredient.Name = strings.TrimSpace(ingredient.Name)
			if ingredient.Name == "" {
				continue
			}
			ingredient.Id = strings.TrimSpace(ingredient.Id)
			if ingredient.Id == "" || taken[ingredient.Id] {
				ingredient.Id = freeIngredientId(ingredient.Name, func(candidate string) bool {
					return taken[candidate]
				})
			}
			if unit, err := cooking.ParseUnit(string(ingredient.DefaultUnit)); err == nil {
				ingredient.DefaultUnit = unit
			} else {
				ingredient.DefaultUnit = cooking.Piece
			}
			ingredient.UpdatedDateTime = now
			catalog = append(catalog, ingredient)
			taken[ingredient.Id] = true
		}

		if err := c.IngredientRepository.ReplaceAll(catalog); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, map[string]int{"count": len(catalog)})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "only GET and POST are supported")
	}
}

// PantryAPI reads and replaces the pantry.
func (c *CookingController) PantryAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pantry, err := c.PantryRepository.List()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		index, err := c.ingredientIndex()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		type pantryLine struct {
			IngredientId string       `json:"ingredientId"`
			Name         string       `json:"name"`
			Amount       float64      `json:"amount"`
			Unit         cooking.Unit `json:"unit"`
		}
		lines := make([]pantryLine, 0, len(pantry))
		for _, item := range pantry {
			lines = append(lines, pantryLine{
				IngredientId: item.IngredientId,
				Name:         ingredientName(index, item.IngredientId),
				Amount:       item.Amount,
				Unit:         item.Unit,
			})
		}
		writeJSON(w, lines)
	case http.MethodPost:
		var incoming []struct {
			IngredientId string  `json:"ingredientId"`
			Name         string  `json:"name"`
			Amount       float64 `json:"amount"`
			Unit         string  `json:"unit"`
		}
		if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		resolver, err := c.newIngredientResolver()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		now := time.Now()
		pantry := []cooking.PantryItem{}
		kept := map[string]bool{}
		for _, line := range incoming {
			label := strings.TrimSpace(line.Name)
			id, err := resolver.resolve(line.IngredientId, label, line.Unit)
			if err != nil {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if id == "" {
				writeJSONError(w, http.StatusBadRequest, "every pantry entry needs an ingredient")
				return
			}

			unit, err := cooking.ParseUnit(line.Unit)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("%s: %s", label, err.Error()))
				return
			}
			if line.Amount < 0 {
				writeJSONError(w, http.StatusBadRequest, label+": the amount must not be negative")
				return
			}
			if kept[id] {
				continue
			}

			pantry = append(pantry, cooking.PantryItem{
				IngredientId:    id,
				Amount:          line.Amount,
				Unit:            unit,
				UpdatedDateTime: now,
			})
			kept[id] = true
		}

		if err := c.PantryRepository.ReplaceAll(pantry); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, map[string]int{"count": len(pantry)})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "only GET and POST are supported")
	}
}

func (c *CookingController) load() ([]cooking.Recipe, map[string]cooking.Ingredient, []cooking.PantryItem, error) {
	recipes, err := c.RecipeRepository.List()
	if err != nil {
		return nil, nil, nil, err
	}
	index, err := c.ingredientIndex()
	if err != nil {
		return nil, nil, nil, err
	}
	pantry, err := c.PantryRepository.List()
	if err != nil {
		return nil, nil, nil, err
	}
	return recipes, index, pantry, nil
}

func (c *CookingController) ingredientIndex() (map[string]cooking.Ingredient, error) {
	ingredients, err := c.IngredientRepository.List()
	if err != nil {
		return nil, err
	}

	index := make(map[string]cooking.Ingredient, len(ingredients))
	for _, ingredient := range ingredients {
		index[ingredient.Id] = ingredient
	}
	return index, nil
}

// savePantry persists the result of a consume operation, dropping the entries
// that were used up.
func (c *CookingController) savePantry(before []cooking.PantryItem, after []cooking.PantryItem) error {
	kept := map[string]bool{}
	for _, item := range after {
		if err := c.PantryRepository.CreateOrUpdate(item); err != nil {
			return err
		}
		kept[item.IngredientId] = true
	}

	for _, item := range before {
		if kept[item.IngredientId] {
			continue
		}
		if err := c.PantryRepository.Delete(item.IngredientId); err != nil {
			return err
		}
	}
	return nil
}

// ingredientResolver maps the free text labels the pages work with onto stable
// catalog ids. Names are matched case insensitively so renaming an ingredient
// in the catalog keeps every recipe and pantry entry pointing at it.
type ingredientResolver struct {
	repository cooking.IIngredientRepository
	byId       map[string]cooking.Ingredient
	byName     map[string]string
}

func (c *CookingController) newIngredientResolver() (*ingredientResolver, error) {
	ingredients, err := c.IngredientRepository.List()
	if err != nil {
		return nil, err
	}

	resolver := &ingredientResolver{
		repository: c.IngredientRepository,
		byId:       make(map[string]cooking.Ingredient, len(ingredients)),
		byName:     make(map[string]string, len(ingredients)),
	}
	for _, ingredient := range ingredients {
		resolver.byId[ingredient.Id] = ingredient
		name := strings.ToLower(strings.TrimSpace(ingredient.Name))
		if _, taken := resolver.byName[name]; name != "" && !taken {
			resolver.byName[name] = ingredient.Id
		}
	}

	return resolver, nil
}

// resolve returns the catalog id for a label, creating the catalog entry when
// the label is new. An empty result means there was nothing to resolve.
func (r *ingredientResolver) resolve(id string, label string, unit string) (string, error) {
	label = strings.TrimSpace(label)
	id = strings.TrimSpace(id)

	if id != "" {
		if known, exists := r.byId[id]; exists && (label == "" || strings.EqualFold(known.Name, label)) {
			return id, nil
		}
	}
	if label == "" {
		if _, exists := r.byId[id]; exists {
			return id, nil
		}
		return "", nil
	}
	if known, exists := r.byName[strings.ToLower(label)]; exists {
		return known, nil
	}

	newId := r.freeId(label)
	ingredient := cooking.Ingredient{
		Id:              newId,
		Name:            label,
		DefaultUnit:     cooking.Piece,
		UpdatedDateTime: time.Now(),
	}
	if parsed, err := cooking.ParseUnit(unit); err == nil {
		ingredient.DefaultUnit = parsed
	}
	if err := r.repository.CreateOrUpdate(ingredient); err != nil {
		return "", err
	}

	r.byId[newId] = ingredient
	r.byName[strings.ToLower(label)] = newId
	return newId, nil
}

// freeId derives an id from a label that no other ingredient is using.
func (r *ingredientResolver) freeId(label string) string {
	return freeIngredientId(label, func(candidate string) bool {
		_, taken := r.byId[candidate]
		return taken
	})
}

func fnv32(value string) uint32 {
	const offset uint32 = 2166136261
	const prime uint32 = 16777619

	hash := offset
	for _, symbol := range value {
		hash ^= uint32(symbol)
		hash *= prime
	}
	return hash
}

// normalizeRecipe validates an incoming recipe, resolves free text ingredient
// names to catalog ids, and renumbers the steps.
func (c *CookingController) normalizeRecipe(recipe *cooking.Recipe) error {
	recipe.Name = strings.TrimSpace(recipe.Name)
	if recipe.Name == "" {
		return fmt.Errorf("the recipe needs a name")
	}
	if recipe.Servings <= 0 {
		recipe.Servings = 1
	}

	resolver, err := c.newIngredientResolver()
	if err != nil {
		return err
	}

	tags := []string{}
	for _, tag := range recipe.Tags {
		if tag = strings.TrimSpace(tag); tag != "" {
			tags = append(tags, tag)
		}
	}
	recipe.Tags = tags

	ingredients, err := normalizeIngredients(resolver, recipe.Ingredients)
	if err != nil {
		return err
	}
	recipe.Ingredients = ingredients

	steps := []cooking.Step{}
	for _, step := range recipe.Steps {
		step.Description = strings.TrimSpace(step.Description)
		if step.Description == "" && len(step.Ingredients) == 0 {
			continue
		}
		if step.Duration < 0 {
			return fmt.Errorf("step %d: the duration must not be negative", len(steps)+1)
		}
		if step.Duration > 0 {
			if step.DurationUnit == "" {
				step.DurationUnit = cooking.Minute
			}
			if !step.DurationUnit.IsValid() {
				return fmt.Errorf("step %d: unknown duration unit %q", len(steps)+1, step.DurationUnit)
			}
		} else {
			step.DurationUnit = ""
		}

		stepIngredients, err := normalizeIngredients(resolver, step.Ingredients)
		if err != nil {
			return err
		}
		step.Ingredients = stepIngredients
		step.Order = len(steps) + 1
		steps = append(steps, step)
	}
	recipe.Steps = steps

	if recipe.Rating < 0 {
		recipe.Rating = 0
	}
	if recipe.Rating > 5 {
		recipe.Rating = 5
	}

	now := time.Now()
	if recipe.Id == "" {
		id, err := c.uniqueRecipeId(recipe.Name)
		if err != nil {
			return err
		}
		recipe.Id = id
		recipe.CreatedDateTime = now
	}
	if recipe.CreatedDateTime.IsZero() {
		recipe.CreatedDateTime = now
	}
	recipe.UpdatedDateTime = now

	return nil
}

func normalizeIngredients(resolver *ingredientResolver, ingredients []cooking.RecipeIngredient) ([]cooking.RecipeIngredient, error) {
	normalized := []cooking.RecipeIngredient{}
	for _, ingredient := range ingredients {
		label := strings.TrimSpace(ingredient.IngredientId)
		if label == "" {
			continue
		}
		if ingredient.Amount < 0 {
			return nil, fmt.Errorf("%s: the amount must not be negative", label)
		}

		unit, err := cooking.ParseUnit(string(ingredient.Unit))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}

		id, err := resolver.resolve("", label, string(unit))
		if err != nil {
			return nil, err
		}
		if id == "" {
			return nil, fmt.Errorf("%s: could not be identified as an ingredient", label)
		}

		ingredient.IngredientId = id
		ingredient.Unit = unit
		ingredient.Note = strings.TrimSpace(ingredient.Note)
		normalized = append(normalized, ingredient)
	}
	return normalized, nil
}

func (c *CookingController) uniqueRecipeId(name string) (string, error) {
	base := slugify(name)
	if base == "" {
		base = "recipe"
	}

	candidate := base
	for suffix := 2; suffix < 1000; suffix++ {
		if _, err := c.RecipeRepository.Get(candidate); err != nil {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, suffix)
	}
	return "", fmt.Errorf("could not generate an id for %q", name)
}

// recipeForEditor rewrites catalog ids back into readable names so the editor
// can work with the labels a person typed.
func recipeForEditor(recipe cooking.Recipe, ingredients []cooking.Ingredient) cooking.Recipe {
	index := make(map[string]cooking.Ingredient, len(ingredients))
	for _, ingredient := range ingredients {
		index[ingredient.Id] = ingredient
	}

	labelled := recipe
	labelled.Ingredients = labelIngredients(recipe.Ingredients, index)
	labelled.Steps = make([]cooking.Step, len(recipe.Steps))
	for i, step := range recipe.OrderedSteps() {
		step.Ingredients = labelIngredients(step.Ingredients, index)
		labelled.Steps[i] = step
	}
	return labelled
}

func labelIngredients(ingredients []cooking.RecipeIngredient, index map[string]cooking.Ingredient) []cooking.RecipeIngredient {
	labelled := make([]cooking.RecipeIngredient, len(ingredients))
	for i, ingredient := range ingredients {
		ingredient.IngredientId = ingredientName(index, ingredient.IngredientId)
		labelled[i] = ingredient
	}
	return labelled
}

func stalenessWeight(recipe cooking.Recipe) int32 {
	if recipe.LastCookedDateTime == nil || recipe.LastCookedDateTime.IsZero() {
		return 30
	}

	days := int32(time.Since(*recipe.LastCookedDateTime).Hours() / 24)
	if days < 1 {
		return 1
	}
	if days > 30 {
		return 30
	}
	return days
}

func matchesQuery(recipe cooking.Recipe, query string) bool {
	if query == "" {
		return true
	}
	if strings.Contains(strings.ToLower(recipe.Name), query) ||
		strings.Contains(strings.ToLower(recipe.Description), query) {
		return true
	}
	for _, tag := range recipe.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func hasTag(recipe cooking.Recipe, tag string) bool {
	if tag == "" {
		return true
	}
	for _, candidate := range recipe.Tags {
		if strings.EqualFold(candidate, tag) {
			return true
		}
	}
	return false
}

func slugify(value string) string {
	var builder strings.Builder
	previousDash := false

	for _, symbol := range strings.ToLower(strings.TrimSpace(value)) {
		switch {
		case unicode.IsLetter(symbol), unicode.IsDigit(symbol):
			builder.WriteRune(symbol)
			previousDash = false
		default:
			if !previousDash && builder.Len() > 0 {
				builder.WriteRune('-')
				previousDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}

// freeIngredientId derives an id from a label that is not already taken.
func freeIngredientId(label string, taken func(string) bool) string {
	base := slugify(label)
	if base == "" {
		// A label made only of symbols, an emoji for instance, still needs a
		// stable key of its own rather than being dropped.
		base = fmt.Sprintf("ingredient-%08x", fnv32(strings.ToLower(strings.TrimSpace(label))))
	}

	candidate := base
	for suffix := 2; ; suffix++ {
		if !taken(candidate) {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, suffix)
	}
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

func render(w http.ResponseWriter, tmpl *template.Template, data any) {
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("[Error] cooking template: %v", err)
	}
}

func redirectWithError(w http.ResponseWriter, r *http.Request, path string, err error) {
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	http.Redirect(w, r, path+separator+"error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("[Error] cooking response: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("[Error] cooking response: %v", err)
	}
}

func mustJSON(payload any) string {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Error] cooking payload: %v", err)
		return "null"
	}
	return string(data)
}
