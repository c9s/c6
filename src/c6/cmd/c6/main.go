package main

import "c6/runtime"
import "github.com/spf13/cobra"
import "fmt"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "c6",
		Short: "C6 is a very fast SASS compatible compiler",
		Long:  `C6 is a SASS compatible implementation written in Go. But wait! this is not only to implement SASS, but also to improve the language for better consistency, syntax and performance.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of C6",
		Long:  `All software has versions. This is C6's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("C6 SASS Compiler v0.1 -- HEAD")
		},
	}
	rootCmd.AddCommand(versionCmd)

	var compileCmd = &cobra.Command{
		Use:   "compile",
		Short: "Compile some scss files",
		// Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run compile!")

			_ = runtime.NewContext()
		},
	}
	rootCmd.AddCommand(compileCmd)
	rootCmd.Execute()
}
