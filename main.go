package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := play(v)
	fmt.Fprint(w, resp)
}

func play(input ArenaUpdate) (response string) {
	log.Printf("IN: %#v", input)

	commands := []string{"F", "R", "L", "T"}

	// Assuming player's id is "self" and opponent's id is "opponent"
	self, ok1 := input.Arena.State["self"]
	opponent, ok2 := input.Arena.State["opponent"]

	if !ok1 || !ok2 {
		// If we don't find either of the players, take a random action
		return commands[rand2.Intn(4)]
	}

	// Check if opponent is in our direct line of sight and close enough
	switch self.Direction {
	case "N":
		if opponent.Y < self.Y && self.X == opponent.X {
			return "T"
		}
	case "S":
		if opponent.Y > self.Y && self.X == opponent.X {
			return "T"
		}
	case "E":
		if opponent.X > self.X && self.Y == opponent.Y {
			return "T"
		}
	case "W":
		if opponent.X < self.X && self.Y == opponent.Y {
			return "T"
		}
	}

	// Check if opponent is to our immediate right or left
	// Only turning, not targeting immediately for simplicity
	if self.Y == opponent.Y {
		if self.X < opponent.X {
			return "R"
		} else {
			return "L"
		}
	} else if self.X == opponent.X {
		if self.Y < opponent.Y {
			return "R"
		} else {
			return "L"
		}
	}

	// If none of the conditions above are met, move forward
	return "F"
}
