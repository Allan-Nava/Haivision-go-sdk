# Haivision Go SDK
[![Go build](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-build.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-build.yml)
[![Go test workflow](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-test.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-test.yml)
[![Release on Tag](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/tag-autorelease.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/tag-autorelease.yml)


The Haivision Go SDK is a software development kit for interacting with Haivision's video streaming platform using the Go programming language. It allows developers to easily integrate Haivision's platform into their Go applications, providing access to a wide range of features such as live streaming, video playback, and media management.

## Installation

To install the Haivision Go SDK, you will need to have Go version 1.13 or later installed on your system. Once Go is installed, you can use the following command to install the SDK:

```bash
go get github.com/Allan-Nava/Haivision-go-sdk
```

### BROADCAST VIDEO ROUTING
Haivision SRT Gateway is a highly flexible and scalable broadcast solution for secure routing of live video streams across different types of IP networks. By serving as a network bridge and converting between protocols including SRT, SRT Gateway provides broadcasters with cost-effective live video streaming to one or multiple destinations for content production and distribution. SRT Gateway includes Haivision’s Path Redundancy feature ensuring uninterrupted live IP video streaming of premium content.

[Rest Api Reference](https://doc.haivision.com/HMG3.7.6/rest-api-integrator-s-reference/rest-api-reference)

## Usage

To use the Haivision Go SDK in your application, you will first need to import it into your code:

```go
import "github.com/Allan-Nava/Haivision-go-sdk"
````

The SDK also provide more functionality such as stop stream, getting stream status, play stream, and more, you can see the full API documentation in the API Reference section.


### Support
If you have any issues or need assistance using the Haivision Go SDK, please contact the developer at allan.nava@hiway.media or visit the project's issue tracker at https://github.com/Allan-Nava/Haivision-go-sdk/issues

### Contribution
We welcome contributions to the Haivision Go SDK. If you would like to contribute, please fork the repository, make your changes, and submit a pull request. When submitting a pull request, please make sure to follow the project's coding style and include tests for any new functionality.

### License
The Haivision Go SDK is released under the MIT License and can be used for both personal and commercial projects.



