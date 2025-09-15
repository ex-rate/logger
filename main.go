package main

import (
	examplesvc "exrate/logger/example_svc"
	"exrate/logger/pkg"
	"log"
)

func main() {
	cfgConsoleOutput := pkg.Config{
		Level:  pkg.DebugLevel,
		Output: pkg.ConsoleOutput,
		Format: "text",
	}

	logger, err := pkg.New(cfgConsoleOutput)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	logger.Info("Logger created successfully")

	serviceLogger := logger.WithService("service1")

	service := examplesvc.NewService(serviceLogger)
	service.Run()
	defer service.Shutdown()

	serviceLogger2 := logger.WithService("service2")

	service2 := examplesvc.NewService(serviceLogger2)
	service2.Run()
	defer service2.Shutdown()

	cfgFileOutputJson := pkg.Config{
		Level:    pkg.DebugLevel,
		Output:   pkg.FileOutput,
		Format:   "json",
		FilePath: "logs/json_output.json",
	}

	logger2, err := pkg.New(cfgFileOutputJson)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	logger2.Info("Logger2 created successfully")

	serviceLogger3 := logger2.WithService("service3")

	service3 := examplesvc.NewService(serviceLogger3)
	service3.Run()
	defer service3.Shutdown()
}
