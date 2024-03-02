
# GoMin Coverage Report 🟢

## Global Coverage
| Statements Covered  | Total Statements | Ratio |
|---|---|---|
|80|553|0.14|true||

## Coverage Errors

| Place | Error |
|---|---|


## Package Coverage
| Package Name | Statements Covered  | Total Statements | Ratio | Valid | Error |
|---|---|---|---|---|---|
|pkg/profiles|3|5|0.60|true||
|pkg/visitor|0|33|0.00|true||
|v0|77|503|0.15|true||
|pkg/declarations|0|12|0.00|true||


## File Coverage

| Package Name | File Name | Statements Covered  | Total Statements | Ratio | Valid | Error |
|---|---|---|---|---|---|---|
|v0|nodes.go|10|10|1.00|true|n/a|
|v0|comparable.go|4|24|0.17|true|n/a|
|v0|writer.go|0|13|0.00|true|n/a|
|v0|commands.go|6|15|0.40|true|n/a|
|v0|construct.go|0|30|0.00|true|n/a|
|v0|parser.go|0|61|0.00|true|n/a|
|v0|matcher.go|5|24|0.21|true|n/a|
|v0|templater.go|16|16|1.00|true|n/a|
|v0|md.go|0|32|0.00|true|n/a|
|v0|reader.go|0|43|0.00|true|n/a|
|v0|processor.go|0|114|0.00|true|n/a|
|v0|ruleset.go|24|24|1.00|true|n/a|
|v0|evaluate.go|12|32|0.38|true|n/a|
|v0|errors.go|0|4|0.00|true|n/a|
|v0|output.go|0|38|0.00|true|n/a|
|v0|format.go|0|23|0.00|true|n/a|
|pkg/declarations|declarations.go|0|12|0.00|true|n/a|
|pkg/profiles|profiles.go|3|5|0.60|true|n/a|
|pkg/visitor|visitor.go|0|33|0.00|true|n/a|


## Block Coverage

| Package Name | File Name | Function Name | Statements Covered  | Total Statements | Ratio | Valid | Error |
|---|---|---|---|---|---|---|---|
|pkg/declarations|declarations.go|New|0|1|0.00|true||
|pkg/declarations|declarations.go|Declarations.DeclByPosition|0|2|0.00|true||
|pkg/declarations|declarations.go|Declarations.search|0|4|0.00|true||
|pkg/declarations|declarations.go|Sort|0|5|0.00|true||
|pkg/profiles|profiles.go|New|3|4|0.75|true||
|pkg/profiles|profiles.go|ProfilesByName.Get|0|1|0.00|true||
|pkg/visitor|visitor.go|Visitor.Visit|0|23|0.00|true||
|pkg/visitor|visitor.go|NewEmplacer|0|1|0.00|true||
|pkg/visitor|visitor.go|FileEmplacer.Emplace|0|5|0.00|true||
|pkg/visitor|visitor.go|NewVisitor|0|4|0.00|true||
|v0|templater.go|NewTemplater|1|1|1.00|true||
|v0|templater.go|NoopReader.Close|1|1|1.00|true||
|v0|templater.go|NoopReader.WriteTo|1|1|1.00|true||
|v0|templater.go|MockExecutor.Execute|1|1|1.00|true||
|v0|templater.go|NoopRenderer.Render|1|1|1.00|true||
|v0|templater.go|Templater.Render|5|5|1.00|true||
|v0|templater.go|NoopReader.Read|1|1|1.00|true||
|v0|templater.go|MockRenderer.Render|1|1|1.00|true||
|v0|templater.go|NewExecutor|4|4|1.00|true||
|v0|md.go|NewMarkdownOutputter|0|4|0.00|true||
|v0|md.go|OSFileFactory.NewFile|0|1|0.00|true||
|v0|md.go|TemplateWriter.Write|0|1|0.00|true||
|v0|md.go|TemplateExecutorFactory.NewExecutor|0|1|0.00|true||
|v0|md.go|NewTemplateWriter|0|1|0.00|true||
|v0|md.go|NewWriterFactory|0|4|0.00|true||
|v0|md.go|NewTemplateExecutorFactory|0|4|0.00|true||
|v0|md.go|TemplateWriterFactory.NewWriter|0|7|0.00|true||
|v0|md.go|NewOSFileFactory|0|4|0.00|true||
|v0|md.go|MarkdownOutputter.WriteOut|0|5|0.00|true||
|v0|reader.go|NewProfileReader|0|1|0.00|true||
|v0|reader.go|ProfileReader.CreateNodeTree|0|42|0.00|true||
|v0|processor.go|FileProcessor.parseFieldIndexListExpr|0|1|0.00|true||
|v0|processor.go|FileProcessor.parseFieldStarExpr|0|1|0.00|true||
|v0|processor.go|FileProcessor.processExprs|0|6|0.00|true||
|v0|processor.go|FileProcessor.parseField|0|1|0.00|true||
|v0|processor.go|FileProcessor.parseCompositeLit|0|1|0.00|true||
|v0|processor.go|FileProcessor.parseStarExpr|0|1|0.00|true||
|v0|processor.go|FileProcessor.parseFieldExpr|0|6|0.00|true||
|v0|processor.go|FileProcessor.parseIndexExpr|0|3|0.00|true||
|v0|processor.go|NewFileProcessor|0|1|0.00|true||
|v0|processor.go|FileProcessor.parseCallExpr|0|3|0.00|true||
|v0|processor.go|FileProcessor.parseFieldIdent|0|1|0.00|true||
|v0|processor.go|FileProcessor.processGenDecl|0|9|0.00|true||
|v0|processor.go|FileProcessor.parseKeyValueExpr|0|3|0.00|true||
|v0|processor.go|FileProcessor.processFuncDecl|0|10|0.00|true||
|v0|processor.go|FileProcessor.processExpr|0|15|0.00|true||
|v0|processor.go|FileProcessor.parseParenExpr|0|1|0.00|true||
|v0|processor.go|FileProcessor.parseUnaryExpr|0|1|0.00|true||
|v0|processor.go|FileProcessor.parseIndexListExpr|0|3|0.00|true||
|v0|processor.go|FileProcessor.processValueSpec|0|10|0.00|true||
|v0|processor.go|FileProcessor.Process|0|12|0.00|true||
|v0|processor.go|FileProcessor.parseFuncLit|0|2|0.00|true||
|v0|processor.go|FileProcessor.parseSliceExpr|0|12|0.00|true||
|v0|processor.go|FileProcessor.parseBinaryExpr|0|3|0.00|true||
|v0|processor.go|FileProcessor.parseFieldIndexExpr|0|1|0.00|true||
|v0|processor.go|FileProcessor.parseTypeAssertExpr|0|6|0.00|true||
|v0|processor.go|FileProcessor.parseEllipsis|0|1|0.00|true||
|v0|ruleset.go|AddMatcher|2|2|1.00|true||
|v0|ruleset.go|AddEvaluator|2|2|1.00|true||
|v0|ruleset.go|NewRuleSet|4|4|1.00|true||
|v0|ruleset.go|ruleSet.match|4|4|1.00|true||
|v0|ruleset.go|ruleSet.Children|7|7|1.00|true||
|v0|ruleset.go|ruleSet.Apply|3|3|1.00|true||
|v0|ruleset.go|AddRuleSet|2|2|1.00|true||
|v0|evaluate.go|eval|12|12|1.00|true||
|v0|evaluate.go|CreateEvaluator|0|5|0.00|true||
|v0|evaluate.go|NewEvaluator|0|1|0.00|true||
|v0|evaluate.go|Evaluator.Evaluate|0|7|0.00|true||
|v0|evaluate.go|ParseOptions|0|7|0.00|true||
|v0|errors.go|CoverageBelowThresholdError.Error|0|1|0.00|true||
|v0|errors.go|InvalidCommandArgumentTypeError.Error|0|1|0.00|true||
|v0|errors.go|InvalidCommandArgumentValueTypeError.Error|0|1|0.00|true||
|v0|errors.go|InvalidMinimumArgumentError.Error|0|1|0.00|true||
|v0|output.go|ErrorFileOutputter.WriteOut|0|8|0.00|true||
|v0|output.go|Open|0|1|0.00|true||
|v0|output.go|NewMultiFileOutputter|0|4|0.00|true||
|v0|output.go|MultiFileOutputter.WriteOut|0|10|0.00|true||
|v0|output.go|NewRecordOutputter|0|4|0.00|true||
|v0|output.go|RecordFileOutputter.WriteOut|0|7|0.00|true||
|v0|output.go|NewErrorOutputter|0|4|0.00|true||
|v0|format.go|createRecord|0|8|0.00|true||
|v0|format.go|nameToString|0|5|0.00|true||
|v0|format.go|Format|0|1|0.00|true||
|v0|format.go|format|0|9|0.00|true||
|v0|matcher.go|NewExactPairMatcher|0|1|0.00|true||
|v0|matcher.go|NewRegexpMatcher|1|1|1.00|true||
|v0|matcher.go|NewNoopMatcher|1|1|1.00|true||
|v0|matcher.go|NewExactMatcher|0|1|0.00|true||
|v0|matcher.go|ExactMatcher.Match|0|6|0.00|true||
|v0|matcher.go|ExactPairMatcher.Match|0|4|0.00|true||
|v0|matcher.go|IndexMatcher.Match|0|4|0.00|true||
|v0|matcher.go|NoopMatcher.Match|0|1|0.00|true||
|v0|matcher.go|NewIndexMatcher|0|1|0.00|true||
|v0|matcher.go|RegexpMatcher.Match|3|4|0.75|true||
|v0|nodes.go|node.Children|1|1|1.00|true||
|v0|nodes.go|node.Leaf|1|1|1.00|true||
|v0|nodes.go|AddNode|2|2|1.00|true||
|v0|nodes.go|AddStatement|2|2|1.00|true||
|v0|nodes.go|NewNode|4|4|1.00|true||
|v0|comparable.go|statementNode.Children|0|1|0.00|true||
|v0|comparable.go|Validate|0|5|0.00|true||
|v0|comparable.go|statements.Previous|0|1|0.00|true||
|v0|comparable.go|ErrorsToRecord|0|4|0.00|true||
|v0|comparable.go|statements.Valid|0|1|0.00|true||
|v0|comparable.go|statements.Covered|1|1|1.00|true||
|v0|comparable.go|ValidatePreviousStatements|0|6|0.00|true||
|v0|comparable.go|evaluatedStatements.Previous|0|1|0.00|true||
|v0|comparable.go|statements.Total|1|1|1.00|true||
|v0|comparable.go|Ratio|1|1|1.00|true||
|v0|comparable.go|evaluatedStatements.Valid|0|1|0.00|true||
|v0|comparable.go|NewStatements|1|1|1.00|true||
|v0|writer.go|Writer.Write|0|4|0.00|true||
|v0|writer.go|NewTabWriterBuilder|0|4|0.00|true||
|v0|writer.go|TabWriterBuilder.NewWriter|0|1|0.00|true||
|v0|writer.go|NewWriter|0|4|0.00|true||
|v0|commands.go|MinimumCommand.Apply|4|4|1.00|true||
|v0|commands.go|NewExcludeCommand|1|1|1.00|true||
|v0|commands.go|ExcludeCommand.Apply|1|1|1.00|true||
|v0|commands.go|NewNoopCommand|0|1|0.00|true||
|v0|commands.go|NoopCommand.Apply|0|1|0.00|true||
|v0|commands.go|NewFallbackCommand|0|1|0.00|true||
|v0|commands.go|FallbackCommand.Apply|0|3|0.00|true||
|v0|commands.go|NewMinimumCommand|0|3|0.00|true||
|v0|construct.go|packageInstanceRuleBuilder.Functions|0|1|0.00|true||
|v0|construct.go|functionsRuleBuilder.genericSurface|0|1|0.00|true||
|v0|construct.go|commandSurface.Parent|0|1|0.00|true||
|v0|construct.go|filesRuleBuilder.genericSurface|0|1|0.00|true||
|v0|construct.go|fileContextRuleBuilder.Functions|0|1|0.00|true||
|v0|construct.go|fileInstanceRuleBuilder.Functions|0|1|0.00|true||
|v0|construct.go|functionsRuleBuilder.Command|0|1|0.00|true||
|v0|construct.go|fileInstanceRuleBuilder.Method|0|1|0.00|true||
|v0|construct.go|AllFiles|0|1|0.00|true||
|v0|construct.go|commandArguments.Value|0|1|0.00|true||
|v0|construct.go|packagesRuleBuilder.Parent|0|1|0.00|true||
|v0|construct.go|filesRuleBuilder.Filter|0|1|0.00|true||
|v0|construct.go|packageInstanceRuleBuilder.Files|0|1|0.00|true||
|v0|construct.go|AllFunctions|0|1|0.00|true||
|v0|construct.go|fileInstanceRuleBuilder.Literal|0|1|0.00|true||
|v0|construct.go|functionsRuleBuilder.Parent|0|1|0.00|true||
|v0|construct.go|functionsRuleBuilder.Filter|0|1|0.00|true||
|v0|construct.go|fileInstanceRuleBuilder.Function|0|1|0.00|true||
|v0|construct.go|filesRuleBuilder.Command|0|1|0.00|true||
|v0|construct.go|packagesRuleBuilder.Command|0|1|0.00|true||
|v0|construct.go|Package|0|1|0.00|true||
|v0|construct.go|packagesRuleBuilder.Filter|0|1|0.00|true||
|v0|construct.go|filesRuleBuilder.Parent|0|1|0.00|true||
|v0|construct.go|commandArguments.Type|0|1|0.00|true||
|v0|construct.go|packageInstanceRuleBuilder.File|0|1|0.00|true||
|v0|construct.go|commandSurface.Command|0|1|0.00|true||
|v0|construct.go|AllPackages|0|1|0.00|true||
|v0|construct.go|packageContextRuleBuilder.Files|0|1|0.00|true||
|v0|construct.go|packageContextRuleBuilder.Functions|0|1|0.00|true||
|v0|construct.go|packagesRuleBuilder.genericSurface|0|1|0.00|true||
|v0|parser.go|parseInt|0|4|0.00|true||
|v0|parser.go|Exclude|0|2|0.00|true||
|v0|parser.go|Min|0|5|0.00|true||
|v0|parser.go|Fallback|0|6|0.00|true||
|v0|parser.go|parseNamePair|0|4|0.00|true||
|v0|parser.go|parseCommandSurface|0|7|0.00|true||
|v0|parser.go|parseCommandArgs|0|22|0.00|true||
|v0|parser.go|parseString|0|4|0.00|true||
|v0|parser.go|parseCommandSurfaces|0|7|0.00|true||

