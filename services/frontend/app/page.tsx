'use client'

import { useState, useRef, useEffect } from 'react'
import { DefaultService } from './client/services/DefaultService'
import { ReportGetResponse } from './client/models/ReportGetResponse'
import { OpenAPI } from './client/core/OpenAPI'

OpenAPI.BASE = process.env.NEXT_PUBLIC_API_BASE ?? 'http://localhost:8080'

export default function Home() {
    const [query, setQuery] = useState('')
    const [reportId, setReportId] = useState<string | null>(null)
    const [report, setReport] = useState<ReportGetResponse | null>(null)
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const pollRef = useRef<NodeJS.Timeout | null>(null)

    const clearPoll = () => {
        if (pollRef.current) {
            clearTimeout(pollRef.current)
            pollRef.current = null
        }
    }

    useEffect(() => {
        return () => clearPoll()
    }, [])

    const pollReport = (id: string) => {
        pollRef.current = setTimeout(async () => {
            try {
                const latest = await DefaultService.getReportById(id)
                if (latest.status === 'ready') {
                    setReport(latest)
                    setLoading(false)
                    clearPoll()
                } else if (latest.status === 'failed') {
                    setError('Report generation failed.')
                    setLoading(false)
                    clearPoll()
                } else {
                    pollReport(id) // still pending
                }
            } catch (e: any) {
                setError(e.message || 'Failed to fetch report.')
                setLoading(false)
                clearPoll()
            }
        }, 1000)
    }

    const generateReport = async () => {
        const trimmed = query.trim()
        if (!trimmed) {
            setError('Please enter a software name to generate a report.')
            setReport(null)
            return
        }

        clearPoll()
        setLoading(true)
        setError(null)
        setReport(null)
        setReportId(null)

        try {
            const created = await DefaultService.createReport({ query: trimmed })
            const id = (created as any).id // adjust if typed differently
            if (!id) throw new Error('No report id returned.')
            setReportId(id)
            pollReport(id)
        } catch (e: any) {
            setError(e.message || 'Failed to start report generation.')
            setLoading(false)
        }
    }

    return (
        <main className="mx-auto max-w-4xl px-5 py-10">
            <div className="space-y-4">
                <div className="flex items-stretch gap-2">
                    <label htmlFor="query" className="sr-only">Software name</label>
                    <div className="flex items-center grow rounded-md bg-white/5 pl-3 outline-1 -outline-offset-1 outline-gray-600 has-[input:focus-within]:outline-2 has-[input:focus-within]:-outline-offset-2 has-[input:focus-within]:outline-indigo-500">
                        <input
                            id="query"
                            name="query"
                            type="text"
                            value={query}
                            onChange={(e) => setQuery(e.target.value)}
                            placeholder="Which software would you like to assess?"
                            className="block min-w-0 grow bg-transparent py-2 pr-3 pl-1 text-base text-white placeholder:text-gray-500 focus:outline-none sm:text-sm/6"
                            disabled={loading}
                        />
                    </div>
                    <button
                        type="button"
                        onClick={generateReport}
                        disabled={loading}
                        className="ml-0 inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-medium text-white hover:bg-indigo-500 disabled:opacity-60"
                    >
                        {loading ? 'Working...' : 'Generate report'}
                    </button>
                </div>

                <div>
                    <div className="rounded-md bg-gray-800 p-4 min-h-[8rem] space-y-2">
                        {error && <p className="text-sm text-red-400">{error}</p>}
                        {!error && loading && !report && (
                            <p className="text-sm text-indigo-300">
                                {reportId ? `Generating report...` : 'Starting report...'}
                            </p>
                        )}
                        {!error && !loading && !report && (
                            <p className="text-sm text-gray-400">Report will appear here after generation.</p>
                        )}
                        {!error && report && (
                            <div className="text-gray-100 space-y-4">
                                <h2 className="text-4xl font-extrabold text-white">
                                    Security report: {report.content?.meta.name}
                                </h2>

                                <div className="space-y-1 text-gray-300">
                                    <p>
                                        <span className="font-semibold text-gray-200">Report ID:</span> {report.id}
                                    </p>
                                    <p>
                                        <span className="font-semibold text-gray-200">Report date:</span>{" "}
                                        {new Date(report.created_at).toLocaleString(undefined, {
                                            year: "numeric",
                                            month: "long",
                                            day: "numeric",
                                            hour: "2-digit",
                                            minute: "2-digit",
                                            second: "2-digit",
                                        })}
                                    </p>
                                </div>

                                <div className="mt-6">
                                    <h3 className="text-2xl font-bold text-white mb-2">Basic Information</h3>
                                    <blockquote className="border-l-4 border-indigo-500 bg-gray-900/50 pl-4 py-2 italic text-gray-200 rounded-md">
                                        {report.content?.meta.short_description}
                                    </blockquote>
                                </div>

                                <div className="space-y-1 text-gray-300 mt-6">
                                    <p>
                                        <span className="font-semibold text-gray-200">Vendor:</span>{" "}
                                        {report.content?.meta.vendor}
                                    </p>
                                    <p>
                                        <span className="font-semibold text-gray-200">Software kind:</span>{" "}
                                        {report.content?.meta.classification}
                                    </p>
                                    <p>
                                        <span className="font-semibold text-gray-200">Alternatives:</span>{" "}
                                        {report.content?.meta.alternatives.join(", ")}
                                    </p>
                                </div>

                                <div className="mt-6">
                                    <h3 className="text-2xl font-bold text-white mb-2">Key Issues</h3>
                                    {report.content?.security_assessment.key_issues.map((issue, index) => (
                                        <div key={index} className="mb-3">
                                            <h4 className="font-semibold text-indigo-400">
                                                {issue.title}
                                            </h4>
                                            <p className="text-gray-300 ">{issue.description}</p>
                                        </div>
                                    ))}
                                </div>

                                <div className="mt-6">
                                    <h3 className="text-2xl font-bold text-white mb-2">
                                        The Verdict
                                    </h3>
                                    <div className="flex items-center gap-4">
                                        {(() => {
                                            const scoreStr = report.content?.security_assessment.security_score ?? '0'
                                            const score = Number.isNaN(Number(scoreStr)) ? 0 : parseInt(String(scoreStr), 10)
                                            const colorClass =
                                                score >= 71 ? 'bg-emerald-500 text-white' :
                                                    score >= 51 ? 'bg-amber-400 text-black' :
                                                        'bg-red-600 text-white'
                                            return (
                                                <>
                                                    <div
                                                        className={`w-20 h-20 rounded-full flex items-center justify-center text-lg font-extrabold ${colorClass}`}
                                                        aria-hidden="true"
                                                    >
                                                        {score}
                                                    </div>
                                                    <div className="text-gray-300">
                                                        <p className="mb-1 font-bold">{report.content?.security_assessment.verdict}</p>
                                                        <p className="text-sm text-gray-400">Security score (0–100)</p>
                                                    </div>
                                                </>
                                            )
                                        })()}
                                    </div>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </main>
    )
}
