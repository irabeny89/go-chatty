package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	goenvy "github.com/irabeny89/go-envy"
	"google.golang.org/genai"
)

func main() {
	goenvy.LoadEnv()
	APIKey := os.Getenv("GEMINI_API_KEY")
	if APIKey == "" {
		log.Fatalln("GEMINI_API_KEY environment variable missing")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Println("Ask Questions")
		cli := bufio.NewScanner(os.Stdin)
		err = cli.Err()
		if err != nil {
			log.Fatal(err)
		}
		cli.Scan()
		input := cli.Text()
		stream := client.Models.GenerateContentStream(
			ctx,
			"gemini-2.5-flash",
			genai.Text(input),
			nil,
		)
		fmt.Println("\n---Answer---")
		for chunk, err := range stream {
			if err != nil {
				log.Fatal(err)
			}
			part := chunk.Candidates[0].Content.Parts[0]
			fmt.Print(part.Text)
		}
		fmt.Println("---End---")
		fmt.Println()
	}
}
