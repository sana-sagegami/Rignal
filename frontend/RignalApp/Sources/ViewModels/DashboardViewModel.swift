import Observation

@Observable
@MainActor
final class DashboardViewModel {
    var summary: SummaryResponse?
    var nextEventTitle: String?
    var nextEventTime: String?
    var isLoading = false
    var errorMessage: String?

    private let backend = BackendService()
    private let calendarService = CalendarService()
    private let notificationService = NotificationService()

    func load(authManager: GoogleAuthManager) async {
        isLoading = true
        errorMessage = nil
        nextEventTitle = nil
        nextEventTime = nil
        defer { isLoading = false }
        do {
            if let user = authManager.currentUser {
                if let event = try await calendarService.fetchNextEvent(user: user) {
                    nextEventTitle = event.title
                    nextEventTime = event.time
                }
            }
            let fetched = try await backend.fetchSummary(nextEvent: nextEventTime)
            summary = fetched
            notificationService.saveWidgetData(
                score: fetched.conditionScore,
                peakStart: fetched.focusPeakStart,
                peakEnd: fetched.focusPeakEnd
            )
            if let bedtime = fetched.recommendBedtime {
                await notificationService.scheduleBedtimeReminder(bedtime: bedtime)
            }
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}
