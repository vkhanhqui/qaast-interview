"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

interface UserLog {
  user_id: string;
  event_type: string;
  details: string;
  created_at: string;
}

export default function AdminUserLogsPage() {
  const [logs, setLogs] = useState<UserLog[]>([]);
  const [cursorStack, setCursorStack] = useState<string[]>([]);
  const [currentCursor, setCurrentCursor] = useState<string | null>(null);
  const [nextCursor, setNextCursor] = useState<string | null>(null);
  const [errorMsg, setErrorMsg] = useState("");
  const [isLastPage, setIsLastPage] = useState(false);

  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;

  async function fetchLogs(cursor?: string) {
    if (!token) return;

    const query = cursor ? `?cursor=${cursor}&limit=10` : "?limit=10";
    const res = await fetch(`http://localhost:8080/admin/userlogs${query}`, {
      headers: { Authorization: `Bearer ${token}` },
    });

    if (!res.ok) {
      setErrorMsg("Failed to load logs");
      return;
    }

    const data = await res.json();
    setLogs(data.user_logs || []);
    setNextCursor(data.next_cursor || null);
  }

  useEffect(() => {
    fetchLogs();
  }, []);

  function handleNext() {
    if (!nextCursor) {
      setErrorMsg("This is the last page.");
      setIsLastPage(true);
      return;
    }
    setCursorStack((prev) => [...prev, currentCursor || ""]);
    setCurrentCursor(nextCursor);
    fetchLogs(nextCursor);
    setErrorMsg("");
  }

  function handlePrev() {
    if (cursorStack.length === 0) return;
    const prevStack = [...cursorStack];
    const prevCursor = prevStack.pop() || null;
    setCursorStack(prevStack);
    setCurrentCursor(prevCursor);
    fetchLogs(prevCursor || undefined);
    setErrorMsg("");
    setIsLastPage(false);
  }

  return (
    <main className="min-h-screen bg-gray-100 p-8 flex flex-col items-center">
      <div className="mb-4 flex items-center gap-4">
        <h1 className="text-2xl font-semibold">User Logs</h1>
        <Link
          href="/admin"
          className="rounded bg-gray-600 px-4 py-2 text-white hover:bg-gray-700"
        >
          Back to Admin Page
        </Link>
      </div>

      {errorMsg && <p className="mb-4 text-red-600">{errorMsg}</p>}

      <div className="overflow-x-auto">
        <table className="mx-auto w-full max-w-5xl border border-gray-300 bg-white text-left shadow">
          <thead className="bg-gray-200">
            <tr>
              <th className="border px-4 py-2">User ID</th>
              <th className="border px-4 py-2">Event Type</th>
              <th className="border px-4 py-2">Details</th>
              <th className="border px-4 py-2">Created At</th>
            </tr>
          </thead>
          <tbody>
            {logs.map((log, idx) => (
              <tr key={`${log.user_id}-${idx}`} className="hover:bg-gray-50">
                <td className="border px-4 py-2 break-all">{log.user_id}</td>
                <td className="border px-4 py-2">{log.event_type}</td>
                <td className="border px-4 py-2 break-all">{log.details}</td>
                <td className="border px-4 py-2">
                  {new Date(log.created_at).toLocaleString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="mt-6 flex gap-4">
        <button
          onClick={handlePrev}
          disabled={cursorStack.length === 0}
          className="rounded bg-gray-500 px-4 py-2 text-white disabled:opacity-40"
        >
          Prev
        </button>
        <button
          onClick={handleNext}
          disabled={isLastPage}
          className="rounded bg-blue-500 px-4 py-2 text-white"
        >
          Next
        </button>
      </div>
    </main>
  );
}
