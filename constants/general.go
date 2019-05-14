package constants

import (
	"time"
)

// GitDirectory ...
const GitDirectory = ".git"

// DockerComposeTokFile ...
const DockerComposeTokFile = "docker-compose.tok.yml"

// HTTPProtocol ...
const HTTPProtocol = "http://"

// HTTPSProtocol ...
const HTTPSProtocol = "https://"

// OneYear is a time.Duration representing a year's worth of seconds.
const OneYear = 8760 * time.Hour

// OneDay is a time.Duration representing a day's worth of seconds.
const OneDay = 24 * time.Hour

// ProjectRootNotFound - Is a flag to mark that the project root was not found when scanning from the users' working directory
const ProjectRootNotFound = "project-root-not-found"
