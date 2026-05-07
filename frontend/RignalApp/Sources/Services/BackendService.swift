import Foundation

struct BackendService {
    private let baseURL: String

    init() {
        guard let url = Bundle.main.infoDictionary?["BackendBaseURL"] as? String else {
            fatalError("BackendBaseURL not set in Info.plist")
        }
        self.baseURL = url
    }

    func triggerAnalysis() async throws {
        guard let url = URL(string: "\(baseURL)/admin/analyze") else { throw URLError(.badURL) }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        let (_, response) = try await URLSession.shared.data(for: request)
        if let http = response as? HTTPURLResponse, !(200..<300).contains(http.statusCode) {
            throw URLError(.badServerResponse)
        }
    }

    func fetchSummary(nextEvent: String? = nil) async throws -> SummaryResponse {
        var components = URLComponents(string: "\(baseURL)/summary")!
        if let nextEvent {
            components.queryItems = [URLQueryItem(name: "next_event", value: nextEvent)]
        }
        guard let url = components.url else { throw URLError(.badURL) }

        let (data, response) = try await URLSession.shared.data(from: url)

        if let http = response as? HTTPURLResponse, !(200..<300).contains(http.statusCode) {
            throw URLError(.badServerResponse)
        }

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        decoder.dateDecodingStrategy = .iso8601
        return try decoder.decode(SummaryResponse.self, from: data)
    }
}
