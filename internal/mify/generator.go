package mify

import (
	"context"
	"fmt"

	"github.com/mify-io/mify/internal/mify/status"
	"github.com/mify-io/mify/internal/mify/util/docker"
	"github.com/mify-io/mify/pkg/generator"
	"github.com/mify-io/mify/pkg/generator/core"
	"github.com/mify-io/mify/pkg/workspace"
)

func ServiceGenerateMany(ctx *CliContext, basePath string, names []string, migrate bool) error {
	descr, err := workspace.InitDescription(basePath)
	if err != nil {
		return err
	}

	if len(names) == 0 {
		names = descr.GetApiServices()
		names = append(names, workspace.DevRunnerName)
	}

	for _, name := range names {
		if err := ServiceGenerate(ctx, basePath, name, migrate); err != nil {
			return fmt.Errorf("service '%s' generation failed: %w", name, err)
		}
	}

	return nil
}

func ServiceGenerate(ctx *CliContext, basePath string, name string, migrate bool) error {
	descr, err := workspace.InitDescription(basePath)
	if err != nil {
		return err
	}

	genPipeline := generator.BuildServicePipeline()
	pg := status.NewGeneratorCliProgress(genPipeline)

	outChan := make(chan core.StepExecResult)

	go genPipeline.Execute(ctx.Ctx, name, descr, migrate, outChan)

	for {
		stepResult := <-outChan
		pg.ReportStep(&stepResult)

		if stepResult.Error != nil {
			return stepResult.Error
		}

		if stepResult.SeqNo == genPipeline.Size()-1 {
			break
		}
	}

	pg.Wait()

	return nil
}

func Cleanup() error {
	err := docker.Cleanup(context.Background())
	if err != nil {
		return err
	}

	return nil
}
