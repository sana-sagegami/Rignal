import GoogleSignIn
import Observation
import UIKit

@Observable
@MainActor
final class GoogleAuthManager {
    var currentUser: GIDGoogleUser?

    var isSignedIn: Bool { currentUser != nil }

    init() {
        currentUser = GIDSignIn.sharedInstance.currentUser
    }

    func restorePreviousSignIn() async {
        do {
            currentUser = try await GIDSignIn.sharedInstance.restorePreviousSignIn()
        } catch {
            currentUser = nil
        }
    }

    func signIn() async throws {
        guard let vc = rootViewController() else { return }
        let result = try await GIDSignIn.sharedInstance.signIn(
            withPresenting: vc,
            hint: nil,
            additionalScopes: ["https://www.googleapis.com/auth/calendar.readonly"]
        )
        currentUser = result.user
    }

    func signOut() {
        GIDSignIn.sharedInstance.signOut()
        currentUser = nil
    }

    private func rootViewController() -> UIViewController? {
        UIApplication.shared.connectedScenes
            .compactMap { $0 as? UIWindowScene }
            .flatMap { $0.windows }
            .first { $0.isKeyWindow }?
            .rootViewController
    }
}
