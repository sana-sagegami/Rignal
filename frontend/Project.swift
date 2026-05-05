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
                ]
            ),
            buildableFolders: [
                "RignalApp/Sources",
                "RignalApp/Resources",
            ],
            dependencies: []
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
