package phpcs

func phpCSXMLTemplate(drupalRoot string) []byte {
	return []byte(`<?xml version="1.0"?>
<ruleset name="DrupalProject">
	<description>Coding standards for a Drupal Project.</description>

	<file>` + drupalRoot + `/sites</file>

	<!--
		Ignore contrib modules and base themes.
	-->
	<exclude-pattern>` + drupalRoot + `/sites/*/libraries/*</exclude-pattern>
	<exclude-pattern>` + drupalRoot + `/sites/*/modules/contrib/*</exclude-pattern>
	<exclude-pattern>` + drupalRoot + `/sites/*/themes/bootstrap/*</exclude-pattern>
	<exclude-pattern>` + drupalRoot + `/sites/*/themes/arena_bootstrap/*</exclude-pattern>

	<!--
		Ignore managed packages and code files.
	-->
	<exclude-pattern>*/vendor/*</exclude-pattern>
	<exclude-pattern>*/node_modules/*</exclude-pattern>

	<!--
		Ignore tests.
	-->
	<exclude-pattern>*/tests/*</exclude-pattern>

	<!--
		Command-line arguments.
		- n Do not print warnings
		- p Show progress of the run
		- s Show sniff codes in all reports
	-->
	<arg value="np"/>

	<!--
		Standards to use.
	-->
	<rule ref="Drupal">
		<exclude name="Drupal.CSS"/>
		<exclude name="Squiz.CSS"/>
	</rule>
	<rule ref="DrupalPractice"/>

	<rule ref="Squiz.WhiteSpace.FunctionOpeningBraceSpace"/>

	<!--
		Ignore long lines in Files.
	-->
	<rule ref="Drupal.Files.LineLength.TooLong">
		<exclude-pattern>*</exclude-pattern>
	</rule>

	<rule ref="Drupal.Files.TxtFileLineLength.TooLong">
		<exclude-pattern>*</exclude-pattern>
	</rule>

	<!--
		Ignore Features generated module files.
	-->
	<exclude-pattern>*/features/*</exclude-pattern>
	<exclude-pattern>*\.features\.*</exclude-pattern>
	<exclude-pattern>*\.strongarm\.*</exclude-pattern>
	<exclude-pattern>*\.context\.*</exclude-pattern>
	<exclude-pattern>*\.field_group\.*</exclude-pattern>
	<exclude-pattern>*\.rules_defaults\.*</exclude-pattern>
	<exclude-pattern>*\.views_default\.*</exclude-pattern>
	<exclude-pattern>*\.linkit_profiles\.*</exclude-pattern>
	<exclude-pattern>*\.file_default_displays\.*</exclude-pattern>
	<exclude-pattern>*\.ds\.*</exclude-pattern>

	<!--
		Some modules use incorrect Naming Conventions.
		-->
	<rule ref="Drupal.NamingConventions.ValidFunctionName.ScopeNotLowerCamel">
		<!-- Context -->
		<exclude-pattern>*_reaction\.*</exclude-pattern>
		<exclude-pattern>*_condition\.*</exclude-pattern>
		<exclude-pattern>*_context_condition_*</exclude-pattern>
		<exclude-pattern>*_context_reaction_*</exclude-pattern>
		<exclude-pattern>*_context_*</exclude-pattern>
		<!-- Entity -->
		<exclude-pattern>*\.entity*</exclude-pattern>
		<!-- CTools -->
		<exclude-pattern>*/plugins/*\.*</exclude-pattern>
	</rule>

	<rule ref="Drupal.NamingConventions.ValidFunctionName.NotLowerCamel">
		<!-- Context -->
		<exclude-pattern>*_context_*</exclude-pattern>
		<!-- CTools -->
		<exclude-pattern>*/plugins/*\.*</exclude-pattern>
	</rule>

	<rule ref="Drupal.NamingConventions.ValidClassName.StartWithCaptial">
		<!-- Context -->
		<exclude-pattern>.*_context_.*</exclude-pattern>
		<!-- CTools -->
		<exclude-pattern>*/plugins/*\.*</exclude-pattern>
	</rule>

	<rule ref="Drupal.NamingConventions.ValidClassName.NoUnderscores">
		<!-- Context -->
		<exclude-pattern>.*_context_.*</exclude-pattern>
		<!-- CTools -->
		<exclude-pattern>*/plugins/*\.*</exclude-pattern>
	</rule>

	<!--
		Ignore global variable Naming Conventions in Drupal settings.
	-->
	<rule ref="Drupal.NamingConventions.ValidGlobal.GlobalUnderScore">
		<!--Settings files-->
		<exclude-pattern>settings\.php</exclude-pattern>
		<exclude-pattern>settings\.*\.php</exclude-pattern>
		<exclude-pattern>*\.settings\.php</exclude-pattern>
	</rule>

	<rule ref="DrupalPractice.General.ClassName.ClassPrefix">
		<!-- Simpletest -->
		<exclude-pattern>*\.test$</exclude-pattern>
		<!-- CTools -->
		<exclude-pattern>*/plugins/*\.*</exclude-pattern>
	</rule>

	<!--
		Views exposed filters forms allowed rewriting of Form State Input.
	-->
	<rule ref="DrupalPractice.General.FormStateInput.Input">
		<exclude-pattern>*\.views\.inc*</exclude-pattern>
	</rule>

	<!--
		Ignore Auto Added Keys for Info Files in Custom modules.
	-->
	<rule ref="Drupal.InfoFiles.AutoAddedKeys">
		<exclude-pattern>*/modules/custom/*.info</exclude-pattern>
	</rule>

</ruleset>
`)
}
