package cooking

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/tthung1997/buddy/core/cooking"
)

func formatAmount(amount float64) string {
	rounded := math.Round(amount*100) / 100
	if rounded == math.Trunc(rounded) {
		return strconv.FormatFloat(rounded, 'f', 0, 64)
	}
	formatted := strconv.FormatFloat(rounded, 'f', 2, 64)
	formatted = strings.TrimRight(formatted, "0")
	return strings.TrimRight(formatted, ".")
}

func formatDuration(duration time.Duration) string {
	if duration <= 0 {
		return ""
	}
	if duration < time.Minute {
		return fmt.Sprintf("%d sec", int(math.Round(duration.Seconds())))
	}

	totalMinutes := int(math.Round(duration.Minutes()))
	days := totalMinutes / (24 * 60)
	hours := (totalMinutes % (24 * 60)) / 60
	minutes := totalMinutes % 60

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d d", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hr", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d min", minutes))
	}

	return strings.Join(parts, " ")
}

func reasonLabel(reason cooking.MissingReason) string {
	switch reason {
	case cooking.MissingReasonInsufficient:
		return "Not enough"
	case cooking.MissingReasonUnitMismatch:
		return "Unit mismatch"
	default:
		return "Not in pantry"
	}
}

func ingredientName(index map[string]cooking.Ingredient, id string) string {
	if ingredient, exists := index[id]; exists && ingredient.Name != "" {
		return ingredient.Name
	}
	return id
}

func buildIngredientLine(ingredient cooking.RecipeIngredient, index map[string]cooking.Ingredient) IngredientLine {
	return IngredientLine{
		IngredientId: ingredient.IngredientId,
		Name:         ingredientName(index, ingredient.IngredientId),
		Amount:       ingredient.Amount,
		AmountLabel:  formatAmount(ingredient.Amount),
		Unit:         ingredient.Unit,
		Optional:     ingredient.Optional,
		Note:         ingredient.Note,
		ImageURL:     index[ingredient.IngredientId].ImageURL,
	}
}

func buildIngredientLines(ingredients []cooking.RecipeIngredient, index map[string]cooking.Ingredient) []IngredientLine {
	lines := make([]IngredientLine, 0, len(ingredients))
	for _, ingredient := range ingredients {
		lines = append(lines, buildIngredientLine(ingredient, index))
	}
	return lines
}

func buildRecipeView(recipe cooking.Recipe, index map[string]cooking.Ingredient) RecipeView {
	steps := recipe.OrderedSteps()
	stepViews := make([]StepView, 0, len(steps))
	for _, step := range steps {
		durationLabel := ""
		if duration, err := step.TotalDuration(); err == nil {
			durationLabel = formatDuration(duration)
		}
		stepViews = append(stepViews, StepView{
			Order:         step.Order,
			Description:   step.Description,
			ImageURL:      step.ImageURL,
			DurationLabel: durationLabel,
			Ingredients:   buildIngredientLines(step.Ingredients, index),
		})
	}

	totalDurationLabel := ""
	if total, err := recipe.TotalDuration(); err == nil {
		totalDurationLabel = formatDuration(total)
	}

	lastCookedLabel := ""
	if recipe.LastCookedDateTime != nil && !recipe.LastCookedDateTime.IsZero() {
		lastCookedLabel = recipe.LastCookedDateTime.Local().Format("Jan 2, 2006")
	}

	ingredients := buildIngredientLines(recipe.RequiredIngredients(), index)

	return RecipeView{
		Id:                 recipe.Id,
		Name:               recipe.Name,
		Description:        recipe.Description,
		ImageURL:           recipe.ImageURL,
		Servings:           recipe.Servings,
		Tags:               recipe.Tags,
		Rating:             recipe.Rating,
		Notes:              recipe.Notes,
		Ingredients:        ingredients,
		Steps:              stepViews,
		TotalDurationLabel: totalDurationLabel,
		LastCookedLabel:    lastCookedLabel,
		StepCount:          len(stepViews),
		IngredientCount:    len(ingredients),
	}
}

func buildMatchView(match cooking.RecipeMatch, index map[string]cooking.Ingredient) MatchView {
	have := make([]IngredientLine, 0, len(match.Have))
	for _, matched := range match.Have {
		have = append(have, IngredientLine{
			IngredientId: matched.IngredientId,
			Name:         ingredientName(index, matched.IngredientId),
			Amount:       matched.Required,
			AmountLabel:  formatAmount(matched.Required),
			Unit:         matched.Unit,
			Optional:     matched.Optional,
			ImageURL:     index[matched.IngredientId].ImageURL,
		})
	}

	missing := make([]MissingLine, 0, len(match.Missing))
	blocking := 0
	for _, item := range match.Missing {
		if !item.Optional {
			blocking++
		}
		missing = append(missing, MissingLine{
			IngredientLine: IngredientLine{
				IngredientId: item.IngredientId,
				Name:         ingredientName(index, item.IngredientId),
				Amount:       item.Amount,
				AmountLabel:  formatAmount(item.Amount),
				Unit:         item.Unit,
				Optional:     item.Optional,
				ImageURL:     index[item.IngredientId].ImageURL,
			},
			Reason:      item.Reason,
			ReasonLabel: reasonLabel(item.Reason),
		})
	}

	statusLabel := "Missing " + strconv.Itoa(blocking) + " ingredient"
	if blocking != 1 {
		statusLabel += "s"
	}
	statusClass := "cooking-status-missing"
	if match.Cookable {
		statusLabel = "Ready to cook"
		statusClass = "cooking-status-ready"
	}

	return MatchView{
		Recipe:          buildRecipeView(match.Recipe, index),
		Servings:        match.Servings,
		Coverage:        match.Coverage,
		CoveragePercent: int(math.Round(match.Coverage * 100)),
		Cookable:        match.Cookable,
		Have:            have,
		Missing:         missing,
		MissingCount:    blocking,
		StatusLabel:     statusLabel,
		StatusClass:     statusClass,
	}
}

func buildMatchViews(matches []cooking.RecipeMatch, index map[string]cooking.Ingredient) []MatchView {
	views := make([]MatchView, 0, len(matches))
	for _, match := range matches {
		views = append(views, buildMatchView(match, index))
	}
	return views
}
