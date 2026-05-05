@preconcurrency import GoogleAPIClientForREST_Calendar
import GoogleSignIn
import Foundation

struct CalendarService {
    func fetchNextEvent(user: GIDGoogleUser) async throws -> (title: String, time: String)? {
        let cal = Calendar.current
        let now = Date()
        let tomorrowStart = cal.startOfDay(for: cal.date(byAdding: .day, value: 1, to: now)!)
        let tomorrowEnd = cal.date(byAdding: .day, value: 1, to: tomorrowStart)!

        let service = GTLRCalendarService()
        service.authorizer = user.fetcherAuthorizer

        let query = GTLRCalendarQuery_EventsList.query(withCalendarId: "primary")
        query.timeMin = GTLRDateTime(date: tomorrowStart)
        query.timeMax = GTLRDateTime(date: tomorrowEnd)
        query.singleEvents = true
        query.orderBy = kGTLRCalendarOrderByStartTime
        query.maxResults = 10

        return try await withCheckedThrowingContinuation { continuation in
            service.executeQuery(query) { _, result, error in
                if let error {
                    continuation.resume(throwing: error)
                    return
                }
                let items = (result as? GTLRCalendar_Events)?.items ?? []
                guard let event = items.first(where: { $0.start?.dateTime != nil }),
                      let date = event.start?.dateTime?.date else {
                    continuation.resume(returning: nil)
                    return
                }
                let f = DateFormatter()
                f.dateFormat = "HH:mm"
                f.timeZone = .current
                continuation.resume(returning: (title: event.summary ?? "予定", time: f.string(from: date)))
            }
        }
    }
}
