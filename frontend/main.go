package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	appCooking "github.com/tthung1997/buddy/app/cooking"
	cookingRepository "github.com/tthung1997/buddy/app/cooking/repository"
	appRandom "github.com/tthung1997/buddy/app/random"
	"github.com/tthung1997/buddy/core/bgg"
	coreRandom "github.com/tthung1997/buddy/core/random"
	"github.com/tthung1997/buddy/frontend/boardgames"
	"github.com/tthung1997/buddy/frontend/cooking"
	"github.com/tthung1997/buddy/frontend/home"
	"github.com/tthung1997/buddy/frontend/shopping"
)

var bggClient *bgg.Client
var cookingController *cooking.CookingController
var randomizer coreRandom.IRandomizer = appRandom.NewSimpleRandomizer()

const cookingDbDir = "frontend/cooking/.db"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("[Info] No .env file found, reading environment variables from system")
	}
	bggToken := os.Getenv("BGG_API_TOKEN")
	if bggToken == "" {
		log.Println("[Warning] BGG_API_TOKEN is not set; BGG API requests may be rejected")
	}
	bggClient = bgg.NewClient(bgg.ClientConfig{
		Root:                bgg.Root,
		BearerToken:         bggToken,
		MaxRetries:          10,
		RetryDelayInSeconds: 5,
	})

	recipes, err := cookingRepository.NewLocalRecipeRepository(cookingDbDir + "/recipes.json")
	if err != nil {
		log.Fatal(err)
	}
	ingredients, err := cookingRepository.NewLocalIngredientRepository(cookingDbDir + "/ingredients.json")
	if err != nil {
		log.Fatal(err)
	}
	pantry, err := cookingRepository.NewLocalPantryRepository(cookingDbDir + "/pantry.json")
	if err != nil {
		log.Fatal(err)
	}
	cookingController = cooking.NewCookingController(
		recipes,
		ingredients,
		pantry,
		appCooking.NewPantryCookEngine(appCooking.NewSimpleUnitConverter()),
		randomizer,
	)
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Incoming] %v %v?%v", r.Method, r.URL.Path, r.URL.RawQuery)
		f(w, r)
	}
}

func main() {
	// no favicon
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// home
	http.HandleFunc("/", logging(home.Index))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	// board games
	bgController := boardgames.NewBoardGamesController(bggClient, randomizer)
	http.HandleFunc("/boardgames", logging(bgController.Index))
	http.HandleFunc("/boardgames/pick", logging(bgController.Pick))
	// backed/preordered tracker
	http.HandleFunc("/boardgames/backed", logging(bgController.Backed))
	http.HandleFunc("/boardgames/backed/add", logging(bgController.BackedAdd))
	http.HandleFunc("/boardgames/backed/receive", logging(bgController.BackedReceive))
	http.HandleFunc("/boardgames/backed/edit", logging(bgController.BackedEdit))

	// shopping
	http.HandleFunc("/shopping", logging(shopping.Index))
	http.HandleFunc("/shopping/inventory", logging(shopping.InventoryHandler))
	http.HandleFunc("/shopping/list", logging(shopping.ShoppingListHandler))

	// cooking
	http.HandleFunc("/cooking", logging(cookingController.Index))
	http.HandleFunc("/cooking/recipes", logging(cookingController.Recipes))
	http.HandleFunc("/cooking/recipes/view", logging(cookingController.Recipe))
	http.HandleFunc("/cooking/recipes/edit", logging(cookingController.Edit))
	http.HandleFunc("/cooking/recipes/save", logging(cookingController.SaveRecipe))
	http.HandleFunc("/cooking/recipes/delete", logging(cookingController.DeleteRecipe))
	http.HandleFunc("/cooking/pantry", logging(cookingController.PantryPage))
	http.HandleFunc("/cooking/api/ingredients", logging(cookingController.IngredientsAPI))
	http.HandleFunc("/cooking/api/pantry", logging(cookingController.PantryAPI))
	// what can I cook right now
	http.HandleFunc("/cooking/cook", logging(cookingController.Cook))
	http.HandleFunc("/cooking/cook/pick", logging(cookingController.Pick))
	http.HandleFunc("/cooking/cook/complete", logging(cookingController.Complete))
	http.HandleFunc("/cooking/shopping/push", logging(cookingController.PushToShopping))

	log.Println("Listening on port 2210")
	err := http.ListenAndServe(":2210", nil)
	if err != nil {
		log.Fatal(err)
	}
}
