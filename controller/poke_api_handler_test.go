package controller

import (
	"catching-pokemons/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonFromPokeApiSuccess(t *testing.T) {
	c := require.New(t)

	pokemon, err := GetPokemonFromPokeApi("1")
	c.NoError(err)
	c.NotEmpty(pokemon)

	body, err := ioutil.ReadFile("samples/poke_api_read.json")
	c.NoError(err)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, pokemon)

}

func TestGetPokemonFromPokeApiSuccessWithMocks(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	body, err := ioutil.ReadFile("samples/poke_api_read.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", "1")

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(200, string(body)))

	pokemon, err := GetPokemonFromPokeApi("1")
	c.NoError(err)
	c.NotEmpty(pokemon)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, pokemon)

}

func TestGetPokemonFromPokeApiInternalServerError(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	body, err := ioutil.ReadFile("samples/poke_api_read.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", "1")

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(500, string(body)))

	_, err = GetPokemonFromPokeApi("1")
	c.NotNil(err)
	c.EqualError(ErrPokeApiFailed, err.Error())

}

func TestGetPokemosFromPokeApiNotFoundError(t *testing.T) {

	c := require.New(t)
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	id := "bulbasaur"

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(404, ""))

	_, err := GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(ErrPokemonNotFound, err.Error())

}

func TestGetPokemon(t *testing.T) {
	c := require.New(t)

	r, err := http.NewRequest("GET", "/pokemon/{id}", nil)
	c.NoError(err)

	w := httptest.NewRecorder()

	//mock path params

	vars := map[string]string{
		"id": "1",
	}

	r = mux.SetURLVars(r, vars)

	GetPokemon(w, r)

	expectedBodyResponse, err := ioutil.ReadFile("samples/api_response.json")
	c.NoError(err)

	var expectedPokemon models.Pokemon

	err = json.Unmarshal([]byte(expectedBodyResponse), &expectedPokemon)
	c.NoError(err)

	var actualPokemon models.Pokemon

	err = json.Unmarshal([]byte(w.Body.Bytes()), &actualPokemon)
	c.NoError(err)

	//check status
	c.Equal(http.StatusOK, w.Code)
	//check body
	c.Equal(expectedPokemon, actualPokemon)
}
