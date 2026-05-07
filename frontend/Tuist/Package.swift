// swift-tools-version: 6.0
import PackageDescription

#if TUIST
    import struct ProjectDescription.PackageSettings

    let packageSettings = PackageSettings(
        // GoogleSignIn/GoogleSignInSwift のみ dynamic（リソースバンドルが必要なため）
        // GoogleAPIClientForREST 系・GTMSessionFetcher 系は static のまま
        // → "GTMSessionFetcherCore is a static product" 警告は出るが無害
        productTypes: [
            "GoogleSignIn": .framework,
            "GoogleSignInSwift": .framework,
            "AppAuth": .framework,
            "GTMAppAuth": .framework,
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
