import ProjectDescription

let project = Project(
    name: "RignalApp",
    targets: [
        .target(
            name: "RignalApp",
            destinations: .iOS,
            product: .app,
            bundleId: "me.bysana.rignal",
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
                    "GIDClientID": "1017409461293-g590kpdc1vnt34sui1j0jua5eag2vk4m.apps.googleusercontent.com",
                    "CFBundleURLTypes": [
                        [
                            "CFBundleTypeRole": "Editor",
                            "CFBundleURLSchemes": ["com.googleusercontent.apps.1017409461293-g590kpdc1vnt34sui1j0jua5eag2vk4m"],
                        ]
                    ],
                ]
            ),
            buildableFolders: [
                "RignalApp/Sources",
                "RignalApp/Resources",
            ],
            entitlements: .dictionary([
                "com.apple.security.application-groups": .array([.string("group.me.bysana.rignal")]),
            ]),
            dependencies: [
                .external(name: "GoogleSignIn"),
                .external(name: "GoogleSignInSwift"),
                .external(name: "AppAuth"),
                .external(name: "GTMAppAuth"),
                .external(name: "GoogleAPIClientForREST_Calendar"),
                .sdk(name: "WidgetKit", type: .framework),
                .target(name: "RignalWidget"),
            ],
            settings: .settings(
                base: ["OTHER_LDFLAGS": "$(inherited) -ObjC"]
            ),
        ),
        .target(
            name: "RignalWidget",
            destinations: .iOS,
            product: .appExtension,
            bundleId: "me.bysana.rignal.RignalWidget",
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
                "com.apple.security.application-groups": .array([.string("group.me.bysana.rignal")]),
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
            bundleId: "me.bysana.rignalTests",
            infoPlist: .default,
            buildableFolders: [
                "RignalApp/Tests"
            ],
            dependencies: [.target(name: "RignalApp")]
        ),
    ]
)
