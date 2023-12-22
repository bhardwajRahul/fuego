package views

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"

	"simple-crud/store"
	"simple-crud/templa"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

// Ressource is the struct that holds useful sources of informations available for the controllers.
type Ressource struct {
	DosingQueries      DosingRepository
	RecipesQueries     RecipeRepository
	IngredientsQueries IngredientRepository
}

func (rs Ressource) showRecipesStd(w http.ResponseWriter, r *http.Request) {
	recipes, err := rs.RecipesQueries.GetRecipes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fp := path.Join("templates", "recipes.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, recipes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rs Ressource) showIndex(c fuego.Ctx[any]) (fuego.Templ, error) {
	recipes, err := rs.RecipesQueries.GetRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	fastRecipes, err := rs.RecipesQueries.GetRandomRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	healthyRecipes, err := rs.RecipesQueries.GetRandomRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	popularRecipes, err := rs.RecipesQueries.GetRandomRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	return templa.Home(templa.HomeProps{
		Recipes:        recipes,
		PopularRecipes: popularRecipes,
		FastRecipes:    fastRecipes,
		HealthyRecipes: healthyRecipes,
	}), nil
}

func (rs Ressource) showRecipes(c fuego.Ctx[any]) (fuego.Templ, error) {
	recipes, err := rs.RecipesQueries.GetRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	return templa.RecipeList(templa.RecipeListProps{
		Recipes: recipes,
	}), nil
}

func (rs Ressource) showSingleRecipes2(c fuego.Ctx[any]) (fuego.Templ, error) {
	id := c.QueryParam("id")

	recipe, err := rs.RecipesQueries.GetRecipe(c.Context(), id)
	if err != nil {
		return nil, fmt.Errorf("error getting recipe %s: %w", id, err)
	}

	ingredients, err := rs.IngredientsQueries.GetIngredientsOfRecipe(c.Context(), id)
	if err != nil {
		slog.Error("Error getting ingredients of recipe", "error", err)
	}

	relatedRecipes, err := rs.RecipesQueries.GetRandomRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	adminCookie, _ := c.Request().Cookie("admin")

	return templa.RecipePage(templa.RecipePageProps{
		Recipe:         recipe,
		Ingredients:    ingredients,
		RelatedRecipes: relatedRecipes,
	}, templa.GeneralProps{
		IsAdmin: adminCookie != nil,
	}), nil
}

func (rs Ressource) searchRecipes(c fuego.Ctx[any]) (fuego.HTML, error) {
	search := c.QueryParam("q")

	recipes, err := rs.RecipesQueries.SearchRecipes(c.Context(), "%"+search+"%")
	if err != nil {
		return "", err
	}

	slog.Debug("recipes", "recipes", recipes, "search", search)

	return c.Render("pages/search.page.html", fuego.H{
		"Recipes": recipes,
		"Search":  search,
		"Filters": fuego.H{
			"Types":       []any{"Entrée", "Plat", "Dessert", "Apéritif"},
			"Ingredients": []any{"Poivron", "Tomate", "Oignon", "Ail", "Piment"},
		},
	})
}

func (rs Ressource) showRecipesList(c fuego.Ctx[any]) (fuego.HTML, error) {
	recipes, err := rs.RecipesQueries.SearchRecipes(c.Context(), "%"+c.QueryParam("search")+"%")
	if err != nil {
		return "", err
	}

	slog.Debug("recipes", "recipes", recipes, "search", c.QueryParam("search"))

	return c.Render("partials/recipes-list.partial.html", recipes)
}

func (rs Ressource) addRecipe(c fuego.Ctx[store.CreateRecipeParams]) (fuego.HTML, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	body.ID = uuid.NewString()

	_, err = rs.RecipesQueries.CreateRecipe(c.Context(), body)
	if err != nil {
		return "", err
	}

	recipes, err := rs.RecipesQueries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/admin.page.html", fuego.H{
		"Recipes": recipes,
	})
}

func (rs Ressource) recipePage(c fuego.Ctx[any]) (fuego.HTML, error) {
	id := c.QueryParam("id")

	recipe, err := rs.RecipesQueries.GetRecipe(c.Context(), id)
	if err != nil {
		return "", fmt.Errorf("error getting recipe %s: %w", id, err)
	}

	ingredients, err := rs.IngredientsQueries.GetIngredientsOfRecipe(c.Context(), id)
	if err != nil {
		slog.Error("Error getting ingredients of recipe", "error", err)
	}

	return c.Render("pages/recipe.page.html", fuego.H{
		"Recipe":       recipe,
		"Ingredients":  ingredients,
		"Instructions": fuego.Markdown(recipe.Instructions),
	})
}

type RecipeRepository interface {
	CreateRecipe(ctx context.Context, arg store.CreateRecipeParams) (store.Recipe, error)
	DeleteRecipe(ctx context.Context, id string) error
	GetRecipe(ctx context.Context, id string) (store.Recipe, error)
	UpdateRecipe(ctx context.Context, arg store.UpdateRecipeParams) (store.Recipe, error)
	GetRecipes(ctx context.Context) ([]store.Recipe, error)
	GetRandomRecipes(ctx context.Context) ([]store.Recipe, error)
	SearchRecipes(ctx context.Context, name string) ([]store.Recipe, error)
}

var _ RecipeRepository = (*store.Queries)(nil)
