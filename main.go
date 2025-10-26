package main

import (
	"bufio"
	"context"
	"fmt"
	"iter"
	"log"
	"os"

	goenvy "github.com/irabeny89/go-envy"
	"google.golang.org/genai"
)

func withStream(ctx context.Context, streamCha chan iter.Seq2[*genai.GenerateContentResponse, error], scanner *bufio.Scanner, client *genai.Client) {
	input := scanner.Text()
	stream := client.Models.GenerateContentStream(
		ctx,
		"gemini-2.5-flash",
		genai.Text(input),
		nil,
	)
	streamCha <- stream
}
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
		fmt.Println("> Ask Questions <")
		fmt.Print("> ")
		cli := bufio.NewScanner(os.Stdin)
		err = cli.Err()
		if err != nil {
			log.Fatal(err)
		}
		cli.Scan()
		streamCha := make(chan iter.Seq2[*genai.GenerateContentResponse, error])
		go withStream(ctx, streamCha, cli, client)
		fmt.Println("\n---Please wait for answer---")
		for chunk, err := range <-streamCha {
			if err != nil {
				log.Fatal(err)
			}
			part := chunk.Candidates[0].Content.Parts[0]
			fmt.Print(part.Text)
		}
		close(streamCha)
		fmt.Println("\n---End---")
		fmt.Println()
	}
}
