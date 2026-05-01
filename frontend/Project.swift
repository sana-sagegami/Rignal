import ProjectDescription

let project = Project(
    name: "AutoZenApp",
    targets: [
        .target(
            name: "AutoZenApp",
            destinations: .iOS,
            product: .app,
            bundleId: "dev.tuist.AutoZenApp",
            infoPlist: .extendingDefault(
                with: [
                    "UILaunchScreen": [
                        "UIColorName": "",
                        "UIImageName": "",
                    ],
                ]
            ),
            buildableFolders: [
                "AutoZenApp/Sources",
                "AutoZenApp/Resources",
            ],
            dependencies: []
        ),
        .target(
            name: "AutoZenAppTests",
            destinations: .iOS,
            product: .unitTests,
            bundleId: "dev.tuist.AutoZenAppTests",
            infoPlist: .default,
            buildableFolders: [
                "AutoZenApp/Tests"
            ],
            dependencies: [.target(name: "AutoZenApp")]
        ),
    ]
)
