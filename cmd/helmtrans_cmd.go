package cmd

import (
    "github.com/spf13/cobra"
    helmtrans "helmtrans/src"
)

var filePath, outPutPath string

//init
func init() {
    //migration data
    helmtransYamltoHelmCmd.Flags().StringVarP(&filePath, "path", "p", "", "file path")
    helmtransYamltoHelmCmd.Flags().StringVarP(&outPutPath, "outpath", "o", "output", "output file path")
    rootCmd.AddCommand(helmtransYamltoHelmCmd)

    helmtransHelmtoYamlCmd.Flags().StringVarP(&filePath, "path", "p", "", "file path")
    rootCmd.AddCommand(helmtransHelmtoYamlCmd)

    helmtransCheckCmd.Flags().StringVarP(&filePath, "path", "p", "", "file path")
    rootCmd.AddCommand(helmtransCheckCmd)
}

//register deploy command
var helmtransYamltoHelmCmd = &cobra.Command{
    Use:   "yamltohelm",
    Short: "Transform yaml to helm",
    Long: `Transform yaml to helm`,
    Run: func(cmd *cobra.Command, args []string) {
        //TO-DO
        helmtrans.YamltoHelm(filePath, outPutPath)
    },
}

//register reset command
var helmtransHelmtoYamlCmd = &cobra.Command{
    Use:   "helmtoyaml",
    Short: "Transform helm to yaml",
    Long: `Transform helm to yaml`,
    Run: func(cmd *cobra.Command, args []string) {
        //TO-DO
    },
}

//register analyse command
var helmtransCheckCmd = &cobra.Command{
    Use:   "check",
    Short: "Check helm",
    Long: `Check helm`,
    Run: func(cmd *cobra.Command, args []string) {
        //TO-DO

    },
}






