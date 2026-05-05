// swift-tools-version: 6.0
import PackageDescription

#if TUIST
    import struct ProjectDescription.PackageSettings

    let packageSettings = PackageSettings(
        productTypes: [
            "GoogleSignIn": .framework,
            "GoogleSignInSwift": .framework,
            "GoogleAPIClientForREST_Calendar": .framework,
            "GTMSessionFetcherCore": .framework,
            "GTMSessionFetcherFull": .framework,
        ]
    )
#endif

let package = Package(
    name: "RignalApp",
    dependencies: [
        .package(url: "https://github.com/google/GoogleSignIn-iOS", from: "7.0.0"),
        .package(url: "https://github.com/google/google-api-objectivec-client-for-rest", from: "3.0.0"),
    ]
)
