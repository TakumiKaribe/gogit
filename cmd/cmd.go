package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var Commands = []*cobra.Command{add, catFile, checkout, commit, hashObject, _init, log, lsTree, merge, rebase, revParse, rm, showRef, tag}

var add = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long:  "add",
	Args:  cobra.MinimumNArgs(1),
	Run:   runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	fmt.Println("add: ", strings.Join(args, " "))
}

var catFile = &cobra.Command{
	Use:   "catFile",
	Short: "catFile",
	Long:  "catFile",
	Args:  cobra.MinimumNArgs(1),
	Run:   runCatFile,
}

func runCatFile(cmd *cobra.Command, args []string) {
	fmt.Println("catFile: ", strings.Join(args, " "))
}

var checkout = &cobra.Command{
	Use:   "checkout",
	Short: "checkout",
	Long:  "checkout",
	Args:  cobra.MinimumNArgs(1),
	Run:   runCheckout,
}

func runCheckout(cmd *cobra.Command, args []string) {
	fmt.Println("checkout: ", strings.Join(args, " "))
}

var commit = &cobra.Command{
	Use:   "commit",
	Short: "commit",
	Long:  "commit",
	Args:  cobra.MinimumNArgs(1),
	Run:   runCommit,
}

func runCommit(cmd *cobra.Command, args []string) {
	fmt.Println("commit: ", strings.Join(args, " "))
}

var hashObject = &cobra.Command{
	Use:   "hashObject",
	Short: "hashObject",
	Long:  "hashObject",
	Args:  cobra.MinimumNArgs(1),
	Run:   runHashObject,
}

func runHashObject(cmd *cobra.Command, args []string) {
	fmt.Println("hashObject: ", strings.Join(args, " "))
}

var _init = &cobra.Command{
	Use:   "init",
	Short: "init",
	Long:  "init",
	Args:  cobra.MinimumNArgs(1),
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("init: ", strings.Join(args, " "))
}

var log = &cobra.Command{
	Use:   "log",
	Short: "log",
	Long:  "log",
	Args:  cobra.MinimumNArgs(1),
	Run:   runLog,
}

func runLog(cmd *cobra.Command, args []string) {
	fmt.Println("log: ", strings.Join(args, " "))
}

var lsTree = &cobra.Command{
	Use:   "lsTree",
	Short: "lsTree",
	Long:  "lsTree",
	Args:  cobra.MinimumNArgs(1),
	Run:   runLsTree,
}

func runLsTree(cmd *cobra.Command, args []string) {
	fmt.Println("lsTree: ", strings.Join(args, " "))
}

var merge = &cobra.Command{
	Use:   "merge",
	Short: "merge",
	Long:  "merge",
	Args:  cobra.MinimumNArgs(1),
	Run:   runMerge,
}

func runMerge(cmd *cobra.Command, args []string) {
	fmt.Println("merge: ", strings.Join(args, " "))
}

var rebase = &cobra.Command{
	Use:   "rebase",
	Short: "rebase",
	Long:  "rebase",
	Args:  cobra.MinimumNArgs(1),
	Run:   runRebase,
}

func runRebase(cmd *cobra.Command, args []string) {
	fmt.Println("rebase: ", strings.Join(args, " "))
}

var revParse = &cobra.Command{
	Use:   "revParse",
	Short: "revParse",
	Long:  "revParse",
	Args:  cobra.MinimumNArgs(1),
	Run:   runRevParse,
}

func runRevParse(cmd *cobra.Command, args []string) {
	fmt.Println("revParse: ", strings.Join(args, " "))
}

var rm = &cobra.Command{
	Use:   "rm",
	Short: "rm",
	Long:  "rm",
	Args:  cobra.MinimumNArgs(1),
	Run:   runAdd,
}

func runRm(cmd *cobra.Command, args []string) {
	fmt.Println("rm: ", strings.Join(args, " "))
}

var showRef = &cobra.Command{
	Use:   "showRef",
	Short: "showRef",
	Long:  "showRef",
	Args:  cobra.MinimumNArgs(1),
	Run:   runShowRef,
}

func runShowRef(cmd *cobra.Command, args []string) {
	fmt.Println("showRef: ", strings.Join(args, " "))
}

var tag = &cobra.Command{
	Use:   "tag",
	Short: "tag",
	Long:  "tag",
	Args:  cobra.MinimumNArgs(1),
	Run:   runTag,
}

func runTag(cmd *cobra.Command, args []string) {
	fmt.Println("tag: ", strings.Join(args, " "))
}
