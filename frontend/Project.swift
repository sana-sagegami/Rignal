import ProjectDescription

let project = Project(
    name: "RignalApp",
    targets: [
        .target(
            name: "RignalApp",
            destinations: .iOS,
            product: .app,
            bundleId: "dev.tuist.RignalApp",
            infoPlist: .extendingDefault(
                with: [
                    "UILaunchScreen": [
                        "UIColorName": "",
                        "UIImageName": "",
                    ],
                    "BackendBaseURL": "http://192.168.1.18:8081",
                    "NSAppTransportSecurity": [
                        "NSAllowsArbitraryLoads": true,
                    ],
                    "GIDClientID": "YOUR_CLIENT_ID.apps.googleusercontent.com",
                    "CFBundleURLTypes": [
                        [
                            "CFBundleTypeRole": "Editor",
                            "CFBundleURLSchemes": ["com.googleusercontent.apps.YOUR_CLIENT_ID"],
                        ]
                    ],
                ]
            ),
            buildableFolders: [
                "RignalApp/Sources",
                "RignalApp/Resources",
            ],
            entitlements: .dictionary([
                "com.apple.security.application-groups": .array([.string("group.dev.tuist.RignalApp")]),
            ]),
            dependencies: [
                .external(name: "GoogleSignIn"),
                .external(name: "GoogleSignInSwift"),
                .external(name: "GoogleAPIClientForREST_Calendar"),
                .sdk(name: "WidgetKit", type: .framework),
                .target(name: "RignalWidget"),
            ]
        ),
        .target(
            name: "RignalWidget",
            destinations: .iOS,
            product: .appExtension,
            bundleId: "dev.tuist.RignalApp.RignalWidget",
            infoPlist: .extendingDefault(
                with: [
                    "NSExtension": [
                        "NSExtensionPointIdentifier": "com.apple.widgetkit-extension",
                    ],
                ]
            ),
            buildableFolders: [
                "RignalWidget/Sources",
            ],
            entitlements: .dictionary([
                "com.apple.security.application-groups": .array([.string("group.dev.tuist.RignalApp")]),
            ]),
            dependencies: [
                .sdk(name: "WidgetKit", type: .framework),
                .sdk(name: "SwiftUI", type: .framework),
            ]
        ),
        .target(
            name: "RignalAppTests",
            destinations: .iOS,
            product: .unitTests,
            bundleId: "dev.tuist.RignalAppTests",
            infoPlist: .default,
            buildableFolders: [
                "RignalApp/Tests"
            ],
            dependencies: [.target(name: "RignalApp")]
        ),
    ]
)
