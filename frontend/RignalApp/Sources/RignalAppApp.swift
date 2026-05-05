import SwiftUI
import GoogleSignIn

@main
struct RignalAppApp: App {
    @State private var authManager = GoogleAuthManager()

    var body: some Scene {
        WindowGroup {
            DashboardView(authManager: authManager)
                .task {
                    await authManager.restorePreviousSignIn()
                    await NotificationService().requestAuthorization()
                }
                .onOpenURL { GIDSignIn.sharedInstance.handle($0) }
        }
    }
}
