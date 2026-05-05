import UserNotifications
import WidgetKit
import Foundation

struct NotificationService {
    private let appGroupID = "group.dev.tuist.RignalApp"
    private let reminderID = "rignal.bedtime.reminder"

    // MARK: - Authorization

    func requestAuthorization() async {
        try? await UNUserNotificationCenter.current()
            .requestAuthorization(options: [.alert, .sound])
    }

    // MARK: - Bedtime Reminder

    func scheduleBedtimeReminder(bedtime: Date) async {
        let center = UNUserNotificationCenter.current()
        let settings = await center.notificationSettings()
        guard settings.authorizationStatus == .authorized else { return }

        center.removePendingNotificationRequests(withIdentifiers: [reminderID])

        guard let reminderTime = Calendar.current.date(byAdding: .minute, value: -30, to: bedtime),
              reminderTime > Date() else { return }

        let content = UNMutableNotificationContent()
        content.title = "就寝リマインド"
        content.body = "推奨就寝時刻まで30分です。そろそろ準備を始めましょう。"
        content.sound = .default

        let components = Calendar.current.dateComponents([.year, .month, .day, .hour, .minute], from: reminderTime)
        let trigger = UNCalendarNotificationTrigger(dateMatching: components, repeats: false)
        try? await center.add(UNNotificationRequest(identifier: reminderID, content: content, trigger: trigger))
    }

    // MARK: - Widget Data (App Group UserDefaults)

    func saveWidgetData(score: Int, peakStart: Date?, peakEnd: Date?) {
        let defaults = UserDefaults(suiteName: appGroupID)
        defaults?.set(score, forKey: "conditionScore")
        if let start = peakStart {
            defaults?.set(start.timeIntervalSince1970, forKey: "focusPeakStart")
        } else {
            defaults?.removeObject(forKey: "focusPeakStart")
        }
        if let end = peakEnd {
            defaults?.set(end.timeIntervalSince1970, forKey: "focusPeakEnd")
        } else {
            defaults?.removeObject(forKey: "focusPeakEnd")
        }
        WidgetCenter.shared.reloadAllTimelines()
    }
}
