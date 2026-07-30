# Graph Report - pp-leads-brasil  (2026-07-30)

## Corpus Check
- 727 files · ~2,230,282 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 8642 nodes · 26790 edges · 63 communities detected
- Extraction: 54% EXTRACTED · 46% INFERRED · 0% AMBIGUOUS · INFERRED: 12342 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]
- [[_COMMUNITY_Community 44|Community 44]]
- [[_COMMUNITY_Community 45|Community 45]]
- [[_COMMUNITY_Community 46|Community 46]]
- [[_COMMUNITY_Community 47|Community 47]]
- [[_COMMUNITY_Community 48|Community 48]]
- [[_COMMUNITY_Community 49|Community 49]]
- [[_COMMUNITY_Community 50|Community 50]]
- [[_COMMUNITY_Community 51|Community 51]]
- [[_COMMUNITY_Community 52|Community 52]]
- [[_COMMUNITY_Community 53|Community 53]]
- [[_COMMUNITY_Community 54|Community 54]]
- [[_COMMUNITY_Community 55|Community 55]]
- [[_COMMUNITY_Community 56|Community 56]]
- [[_COMMUNITY_Community 60|Community 60]]
- [[_COMMUNITY_Community 61|Community 61]]
- [[_COMMUNITY_Community 74|Community 74]]
- [[_COMMUNITY_Community 75|Community 75]]
- [[_COMMUNITY_Community 76|Community 76]]
- [[_COMMUNITY_Community 77|Community 77]]

## God Nodes (most connected - your core abstractions)
1. `contains()` - 1419 edges
2. `New()` - 583 edges
3. `writeFile()` - 539 edges
4. `Run()` - 506 edges
5. `minimalSpec()` - 341 edges
6. `Parse()` - 274 edges
7. `CLI()` - 228 edges
8. `Execute()` - 164 edges
9. `runGoCommand()` - 138 edges
10. `Split()` - 100 edges

## Surprising Connections (you probably didn't know these)
- `Load()` --calls--> `TestLoad_FileAbsentReturnsEmpty()`  [INFERRED]
  leads-brasil-pp-cli/internal/config/config.go → cli-printing-press/internal/mcpoverrides/overrides_test.go
- `newRootCmd()` --calls--> `newLearnConfig()`  [INFERRED]
  leads-brasil-pp-cli/internal/cli/root.go → cli-printing-press/testdata/golden/expected/generate-learn-loop-api/learn-loop-example/internal/cli/learn_init.go
- `TestCatalogFSContainsYAMLFiles()` --calls--> `Name()`  [INFERRED]
  cli-printing-press/internal/catalog/catalog_test.go → internal/usecase/config.go
- `_run_telegram()` --calls--> `Run()`  [INFERRED]
  organizejr-pp-leads/telegram_people.py → cli-printing-press/internal/llm/llm.go
- `_telegram_candidates()` --calls--> `Execute()`  [INFERRED]
  organizejr-pp-leads/telegram_people.py → leads-brasil-pp-cli/internal/cli/root.go

## Communities

### Community 0 - "Community 0"
Cohesion: 0.01
Nodes (807): hasOverrideMarker(), renderAgentcookieManifest(), syncKeysFromAuth(), TestWriteAgentcookieManifest_BearerToken(), TestWriteAgentcookieManifest_EnvVarSpecsRespectsSensitivity(), TestWriteAgentcookieManifest_Idempotent(), TestWriteAgentcookieManifest_RespectsOverrideMarker(), TestWriteAgentcookieManifest_SkipsCookieOnly() (+799 more)

### Community 1 - "Community 1"
Cohesion: 0.01
Nodes (559): newAPICmd(), TestAuthCmdHasDoctorSubcommand(), TestAuthDoctorCmdJSONFlag(), backupDirForReplace(), backupFileForReplace(), captureDomain(), commonCaptureBaseURL(), newBrowserSniffCmd() (+551 more)

### Community 2 - "Community 2"
Cohesion: 0.01
Nodes (525): TestAnalysisWarning_StringV2Compat(), TestAuthAnalysis_CandidateTypesV2Compat(), TestProtectionObservation_NotesV2Compat(), TestReadTrafficAnalysis_V2ShapeFile(), TestTrafficAnalysis_GenerationHintsMapCompat(), TestTrafficAnalysis_VersionNormalization(), BearerRefresh, Entry (+517 more)

### Community 3 - "Community 3"
Cohesion: 0.01
Nodes (479): outputIsTerminal(), getEnv, EmbossDelta, EmbossReport, EmbossSnapshot, FeedbackEntry, cliExecutableName(), SiblingCLIPath() (+471 more)

### Community 4 - "Community 4"
Cohesion: 0.01
Nodes (379): clearCollidingParents(), flattenCollidingBodyFields(), TestFlattenCollidingBodyFields_NestedPrefixShape(), TestFlattenCollidingBodyFields_NoCollisionPassesThrough(), countBodyLeaves(), FlattenCollidingBodyFields(), Ident(), deepBodyFixture() (+371 more)

### Community 5 - "Community 5"
Cohesion: 0.01
Nodes (287): WriteTrafficAnalysis(), chromeDataDir(), clearPendingDeviceCode(), cookieToolSupportsProfiles(), copyFileIfExists(), countCookiesForDomain(), detectCookieTool(), discoverChromeProfiles() (+279 more)

### Community 6 - "Community 6"
Cohesion: 0.01
Nodes (421): AnalyzeTraffic(), anyHeaderPrefix(), apiHostsForProtection(), appendEvidence(), bucketConfidence(), buildEndpointClusters(), buildTrafficSummary(), captchaChallengeText() (+413 more)

### Community 7 - "Community 7"
Cohesion: 0.01
Nodes (387): absOrSame(), advertisedNovelFeaturePath(), appendedStringValues(), authCandidatesForPrefix(), authFormatInlineMapPreservesToken(), buildDogfoodBinary(), checkAuth(), checkCommandTree() (+379 more)

### Community 8 - "Community 8"
Cohesion: 0.01
Nodes (377): CleanupOptions, CleanupGeneratedCLI(), removeDirIfExists(), removeFileIfExists(), removeFinderMetadata(), TestCleanupGeneratedCLI(), writeArtifactFile(), Name() (+369 more)

### Community 9 - "Community 9"
Cohesion: 0.01
Nodes (358): addCommands(), addImports(), addPersistentFlags(), addPostExecuteFlush(), addPreRunBlocks(), addRootFlagsFields(), appendExecuteStatementsAfterLast(), checkRootShape() (+350 more)

### Community 10 - "Community 10"
Cohesion: 0.02
Nodes (266): Store, isSpace(), apply_icp_dropdown(), attempt_id_for(), attempt_identity_payload(), attempt_overlaps(), best_contact_channel(), blocking_attempts_for() (+258 more)

### Community 11 - "Community 11"
Cohesion: 0.02
Nodes (231): autoRefreshIfStale(), cachePolicy(), emitCacheRefreshFailedEvent(), ensureFreshForCommand(), ensureFreshForResources(), refreshTimeout(), runAutoRefresh(), accessWarning (+223 more)

### Community 12 - "Community 12"
Cohesion: 0.02
Nodes (266): childIsSchema(), convertToStringKeyed(), decodeSpecTree(), isNameKeyedParent(), normalizeExamplesValue(), normalizeSpecData(), normalizeSpecDataWithMetadata(), normalizeSpecTree() (+258 more)

### Community 13 - "Community 13"
Cohesion: 0.02
Nodes (198): Profile, ColumnDef, IndexDef, TableDef, SyncHint, VisionCustomization, ShardedSubResourceTableName(), SubResourceShardedNames() (+190 more)

### Community 14 - "Community 14"
Cohesion: 0.02
Nodes (190): PIIAllAcceptedIssue, PIIAuditOptions, PIIAuditResult, PIICompletionStatus, piiDetector, PIIFinding, PIILedger, PIILedgerDelta (+182 more)

### Community 15 - "Community 15"
Cohesion: 0.02
Nodes (165): Apply(), assertGitClean(), injectAddCommands(), jaccard(), nonPlaceholderTokens(), parseStmtViaDST(), readModulePaths(), setOf() (+157 more)

### Community 16 - "Community 16"
Cohesion: 0.02
Nodes (171): applyFix(), classifyFailures(), detectHTTPMethod(), findCommandFile(), hasPrintDryRunHelper(), revertFiles(), runBuildChecks(), RunFixLoop() (+163 more)

### Community 17 - "Community 17"
Cohesion: 0.03
Nodes (124): absFixturePath(), assertFileExists(), snapshotTree(), stageFixture(), TestApplyEbayAuthFixturePreservesAuth(), TestApplyIdempotency(), TestApplyPostmanExploreFixture(), dependentPathParamDef (+116 more)

### Community 18 - "Community 18"
Cohesion: 0.04
Nodes (113): ExecutablePath(), ExecutablePathForGOOS(), TestExecutablePathForGOOS(), validationGate, GoRunArgs(), TestGoRunArgsUsesDefaultMode(), commonCommandVerb(), completePartialRedactionSentinel() (+105 more)

### Community 19 - "Community 19"
Cohesion: 0.05
Nodes (63): CasaDadosClient, Client, LeadRecord, cleanEmail(), cnpjFromInput(), collectPhones(), companyDomain(), CompanyPayload() (+55 more)

### Community 20 - "Community 20"
Cohesion: 0.04
Nodes (91): authEnvHintComment(), RecipeIntent, RecipeIntentParam, RecipeIntentParamType, TestMCPDescription(), ParamDescriptionCompactor, paramDescriptionKey, CatalogDescription() (+83 more)

### Community 21 - "Community 21"
Cohesion: 0.05
Nodes (91): autoBundleForHost(), buildBundleBinaries(), buildMCPBBinary(), bundleBinaryArchiveName(), bundleBinaryPath(), newBundleCmd(), resolvePlatform(), installFakeGo() (+83 more)

### Community 22 - "Community 22"
Cohesion: 0.04
Nodes (83): ManifestAuth, manifestBodyParam, manifestEndpointRecord, ManifestHeader, ManifestMCP, ManifestParam, ManifestTier, ManifestTiers (+75 more)

### Community 23 - "Community 23"
Cohesion: 0.05
Nodes (74): Aggregate(), AggregateAuth(), deduplicateStrings(), FilterByHost(), hostFromURL(), NormalizePath(), paramFieldCount(), sortedParams() (+66 more)

### Community 24 - "Community 24"
Cohesion: 0.05
Nodes (61): PlanCommand, planGoModData, planParentCommand, planRootData, PlanSpec, planStreamingData, planSubCommand, parseCopyrightOwner() (+53 more)

### Community 25 - "Community 25"
Cohesion: 0.09
Nodes (56): capturedKeysIndex, SyncParamDropFinding, SyncParamDropResult, callPassedKeys(), canonicalSyncPath(), CheckSyncParamDrop(), collectSyncSourceFiles(), extractCompositeLiteralKeys() (+48 more)

### Community 26 - "Community 26"
Cohesion: 0.08
Nodes (53): canonicalMCPSurfacePath(), cloneStringBoolMap(), cobraConstructorCallName(), cobratreeCommandKind(), codeOrchEndpointCount(), collectCobraSourceCommands(), estimateCobratreeCommandTool(), estimateCobratreeRuntimeTokens() (+45 more)

### Community 27 - "Community 27"
Cohesion: 0.07
Nodes (43): classifyResponse(), hasHeaderPrefix(), isClear(), lowerHeaders(), protectionsToEvidence(), TestClassifyResponse(), TestIsClear(), classifyFailure() (+35 more)

### Community 28 - "Community 28"
Cohesion: 0.09
Nodes (48): authNarrativeMentionsDoctor(), existingReadmeIntroLead(), findLineWithPrefix(), findMarkdownHeading(), findMarkdownHeadingInRange(), findNextLevelTwoHeading(), findNextMarkdownHeadingAtMost(), findReadmeIntroEnd() (+40 more)

### Community 29 - "Community 29"
Cohesion: 0.1
Nodes (47): gqlArg, gqlField, gqlType, addScalarTypes(), addSupportTypes(), buildGetEndpoint(), buildListEndpoint(), buildMutationEndpoint() (+39 more)

### Community 30 - "Community 30"
Cohesion: 0.11
Nodes (44): phase5AcceptancePassed(), phase5Level(), phase5SkipAllowed(), phase5SkipReason(), TestValidatePhase5Gate_APIKeyMissingSkipAllowed(), TestValidatePhase5Gate_CookieAuthNotSkippedByMissingAPIKey(), TestValidatePhase5Gate_CredentialAuthRejectsLANSkipReason(), TestValidatePhase5Gate_FullPassAcceptsAccountedSkips() (+36 more)

### Community 31 - "Community 31"
Cohesion: 0.09
Nodes (38): DetectAPIType(), readHead(), TestDetectAPITypeGraphQL(), TestDetectAPITypeGRPC(), TestDetectAPITypeREST(), TestDetectAPITypeURL(), TestDetectAPITypeYAMLFallback(), writeTempSpec() (+30 more)

### Community 32 - "Community 32"
Cohesion: 0.09
Nodes (37): appendMethodMarker(), collectParams(), Compose(), composeAction(), composeActionWithFallback(), composeOptional(), composeRequired(), composeReturns() (+29 more)

### Community 33 - "Community 33"
Cohesion: 0.09
Nodes (34): applyExamples(), applyHelpTexts(), applyREADME(), escapeGoString(), commandInfo, ExampleSet, HelpImprovement, PolishRequest (+26 more)

### Community 34 - "Community 34"
Cohesion: 0.17
Nodes (20): deriveAuthVerifyPath(), deriveHealthCheckPath(), findMeShapedEndpointPath(), isMeShapedEndpoint(), TestDeriveAuthVerifyPath_FallsBackToEmpty(), TestDeriveAuthVerifyPath_NilSpec(), TestDeriveAuthVerifyPath_PicksMeShapedEndpoint(), TestDeriveAuthVerifyPath_PrioritizesExplicitOverride() (+12 more)

### Community 35 - "Community 35"
Cohesion: 0.17
Nodes (18): canonicalFinding, pythonReport, verifySkillRunResult, checkVerifySkill(), emitMergedJSON(), ExtractInstallSectionForTest(), indentLines(), TestIsWindowsStorePython() (+10 more)

### Community 36 - "Community 36"
Cohesion: 0.14
Nodes (16): detectCollisions(), featureAlreadyColliding(), hasTestSuffix(), looksLikePR218Dropin(), extractOwner(), planDropins(), renderDropin(), cliProvenance (+8 more)

### Community 37 - "Community 37"
Cohesion: 0.16
Nodes (14): newAuthCmd(), newAuthDoctorCmd(), Finding, Status, Summary, orDash(), RenderJSON(), RenderTable() (+6 more)

### Community 38 - "Community 38"
Cohesion: 0.2
Nodes (17): internalSkillFinding, internalSkillFrontmatter, internalSkillReport, hasH1Heading(), newVerifyInternalSkillCmd(), runVerifyInternalSkillChecks(), splitFrontmatter(), TestHasH1Heading() (+9 more)

### Community 39 - "Community 39"
Cohesion: 0.24
Nodes (15): MCPAuditFinding, auditLibraryCLI(), newMCPAuditCmd(), recommendForFinding(), renderMCPAuditTable(), runMCPAudit(), makeMCPTools(), mustWrite() (+7 more)

### Community 40 - "Community 40"
Cohesion: 0.31
Nodes (12): cliArgsFromMCP(), RunCLICommand(), shellOutToCLI(), SplitShellArgs(), TestArgsFieldRejectsFlagLikeTokens(), TestCliArgsFromMCP_AllowsPerCommandFlags(), TestCliArgsFromMCP_BlocksRootFlags(), TestRunCLICommandFallsBackToStdoutOnFailureWithoutStderr() (+4 more)

### Community 41 - "Community 41"
Cohesion: 0.26
Nodes (15): decodeOpenAPIRaw(), detectNestedDataEnvelopeFixtures(), detectNestedDataEnvelopeFixturesFromRaw(), isDogfoodHTTPMethod(), isJSONMediaType(), isRawArraySchema(), nestedDataArrayKey(), resolveRawPointer() (+7 more)

### Community 42 - "Community 42"
Cohesion: 0.2
Nodes (14): ApplyCatalogAuthEnvVars(), ApplyRuntimeMetadata(), IsReplaceableBaseURL(), mergeAuthEnvVars(), RebaseAuthEnvPrefix(), TestApplyCatalogAuthEnvVars_DedupesBetweenCatalogAndExisting(), TestApplyCatalogAuthEnvVars_NoopForBasicAuth(), TestApplyCatalogAuthEnvVars_NoopWhenAuthTypeNone() (+6 more)

### Community 43 - "Community 43"
Cohesion: 0.32
Nodes (11): buildAgentContext(), buildAgentDiscoveryContext(), collectAgentCommands(), newAgentContextCmd(), agentContext, agentContextAuth, agentContextAuthEnvVar, agentContextCLI (+3 more)

### Community 44 - "Community 44"
Cohesion: 0.3
Nodes (7): cookieJar, persistedCookie, cookieJarPath(), LoadCookieJar(), mergeAndWriteCookieRows(), sanitizeCookieValue(), WriteCookieJarFromMap()

### Community 45 - "Community 45"
Cohesion: 0.17
Nodes (11): AuthCapture, EnrichedCapture, EnrichedEntry, HAR, HAREntry, HARHeader, HARLog, HARPostData (+3 more)

### Community 46 - "Community 46"
Cohesion: 0.24
Nodes (8): ApisGuruPattern(), catalogSpecSource(), DiscoverSpec(), TestApisGuruPattern(), TestDiscoverSpec_CatalogEntryWins(), TestDiscoverSpec_KnownAPI(), TestDiscoverSpec_UnknownAPI(), KnownSpec

### Community 47 - "Community 47"
Cohesion: 0.26
Nodes (10): compareGapsAndAdvantages(), RunComparative(), scoreAlternative(), TestCompareGapsAndAdvantages(), TestRunComparative(), TestRunComparativeLoadsResearchFromSiblingResearchDir(), TestScoreAlternative(), writeComparativeReport() (+2 more)

### Community 48 - "Community 48"
Cohesion: 0.31
Nodes (8): DetectAsyncJobs(), detectOne(), findStatusSibling(), responseJobIDField(), TestDetectAsyncJobs(), TestDetectAsyncJobs_NilSpec(), TestResponseJobIDField_VariousNames(), AsyncJobInfo

### Community 49 - "Community 49"
Cohesion: 0.39
Nodes (8): MapEntities(), mapResource(), matchesKeywords(), primaryKeywordsForArchetype(), scanEntityFields(), scanMappingFields(), EntityMapping, WorkflowTemplateContext

### Community 50 - "Community 50"
Cohesion: 0.39
Nodes (6): applyParamPatches(), MergeOverlay(), TestMergeOverlay(), TestMergeOverlayNilSafe(), TestMergeOverlayPublicParamNames(), TestMergeOverlayRejectsEmptyFlagName()

### Community 51 - "Community 51"
Cohesion: 0.33
Nodes (5): CreateProjectRequest, Currency, Project, Task, UpdateTaskRequest

### Community 52 - "Community 52"
Cohesion: 0.53
Nodes (4): deleteKey(), keychainService(), loadKey(), saveKey()

### Community 53 - "Community 53"
Cohesion: 0.33
Nodes (5): AggregatedEndpoint, DiscoveredAuth, DiscoveredEndpoint, DiscoveredParam, SourceResult

### Community 54 - "Community 54"
Cohesion: 0.4
Nodes (4): EndpointOverlay, ParamPatch, ResourceOverlay, SpecOverlay

### Community 55 - "Community 55"
Cohesion: 0.5
Nodes (1): syncLive()

### Community 56 - "Community 56"
Cohesion: 0.83
Nodes (3): findRepoRoot(), TestVerifySkillDriftWorkflowGuardsLibraryCopy(), TestVerifySkillScriptInSync()

### Community 60 - "Community 60"
Cohesion: 1.0
Nodes (1): CafeMenu

### Community 61 - "Community 61"
Cohesion: 1.0
Nodes (1): NamingRule

### Community 74 - "Community 74"
Cohesion: 1.0
Nodes (1): Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_

### Community 75 - "Community 75"
Cohesion: 1.0
Nodes (1): Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_

### Community 76 - "Community 76"
Cohesion: 1.0
Nodes (1): Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_

### Community 77 - "Community 77"
Cohesion: 1.0
Nodes (1): Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_

## Knowledge Gaps
- **513 isolated node(s):** `Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_`, `tokenResponse`, `JobRow`, `WaitOptions`, `CafeMenu` (+508 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 55`** (4 nodes): `sync.go`, `sync.go`, `sync.go`, `syncLive()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 60`** (2 nodes): `types.go`, `CafeMenu`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 61`** (2 nodes): `naming_rules.go`, `NamingRule`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 74`** (1 nodes): `Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 75`** (1 nodes): `Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 76`** (1 nodes): `Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 77`** (1 nodes): `Returns (score_0_100, need_enrichment_YN, problems_str, need_enrichment_reasons_`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `contains()` connect `Community 0` to `Community 1`, `Community 2`, `Community 3`, `Community 4`, `Community 5`, `Community 6`, `Community 7`, `Community 8`, `Community 9`, `Community 10`, `Community 11`, `Community 12`, `Community 13`, `Community 14`, `Community 15`, `Community 16`, `Community 17`, `Community 18`, `Community 19`, `Community 20`, `Community 21`, `Community 22`, `Community 23`, `Community 24`, `Community 25`, `Community 26`, `Community 27`, `Community 28`, `Community 29`, `Community 30`, `Community 31`, `Community 32`, `Community 33`, `Community 34`, `Community 35`, `Community 36`, `Community 37`, `Community 38`, `Community 39`, `Community 40`, `Community 41`, `Community 42`, `Community 46`, `Community 48`, `Community 49`, `Community 50`, `Community 56`?**
  _High betweenness centrality (0.362) - this node is a cross-community bridge._
- **Why does `Run()` connect `Community 2` to `Community 0`, `Community 1`, `Community 3`, `Community 4`, `Community 5`, `Community 6`, `Community 7`, `Community 8`, `Community 9`, `Community 10`, `Community 11`, `Community 12`, `Community 13`, `Community 14`, `Community 15`, `Community 16`, `Community 18`, `Community 19`, `Community 20`, `Community 21`, `Community 22`, `Community 23`, `Community 24`, `Community 26`, `Community 27`, `Community 31`, `Community 32`, `Community 33`, `Community 34`, `Community 35`, `Community 36`, `Community 38`, `Community 40`, `Community 47`, `Community 48`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Why does `writeFile()` connect `Community 1` to `Community 0`, `Community 2`, `Community 3`, `Community 4`, `Community 5`, `Community 6`, `Community 7`, `Community 8`, `Community 9`, `Community 11`, `Community 13`, `Community 14`, `Community 15`, `Community 16`, `Community 17`, `Community 18`, `Community 19`, `Community 20`, `Community 21`, `Community 22`, `Community 24`, `Community 25`, `Community 26`, `Community 28`, `Community 30`, `Community 31`, `Community 33`, `Community 36`, `Community 38`, `Community 39`, `Community 40`, `Community 44`, `Community 47`?**
  _High betweenness centrality (0.095) - this node is a cross-community bridge._
- **Are the 1415 inferred relationships involving `contains()` (e.g. with `main()` and `applyAuthFormat()`) actually correct?**
  _`contains()` has 1415 INFERRED edges - model-reasoned connections that need verification._
- **Are the 574 inferred relationships involving `New()` (e.g. with `NewAdaptiveLimiter()` and `.Set()`) actually correct?**
  _`New()` has 574 INFERRED edges - model-reasoned connections that need verification._
- **Are the 513 inferred relationships involving `writeFile()` (e.g. with `.save()` and `.writeCache()`) actually correct?**
  _`writeFile()` has 513 INFERRED edges - model-reasoned connections that need verification._
- **Are the 504 inferred relationships involving `Run()` (e.g. with `_run_telegram()` and `TestIsCobraUsageError()`) actually correct?**
  _`Run()` has 504 INFERRED edges - model-reasoned connections that need verification._