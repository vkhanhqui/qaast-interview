"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
}

export default function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [cursorStack, setCursorStack] = useState<string[]>([]);
  const [currentCursor, setCurrentCursor] = useState<string | null>(null);
  const [nextCursor, setNextCursor] = useState<string | null>(null);
  const [errorMsg, setErrorMsg] = useState("");
  const [isLastPage, setIsLastPage] = useState(false);

  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;

  async function fetchUsers(cursor?: string) {
    if (!token) return;

    const query = cursor ? `?cursor=${cursor}&limit=10` : "?limit=10";
    const res = await fetch(`http://localhost:8080/admin/users${query}`, {
      headers: { Authorization: `Bearer ${token}` },
    });

    if (!res.ok) {
      setErrorMsg("Failed to load users");
      return;
    }

    const data = await res.json();
    if (!data.users || !data.next_cursor) {
      setIsLastPage(true);
      setErrorMsg("This is the last page.");
      setUsers(data.users || []);
      return;
    }
    setUsers(data.users || []);
    setNextCursor(data.next_cursor || null);
  }

  useEffect(() => {
    fetchUsers();
  }, []);

  function handleNext() {
    if (!nextCursor) {
      setErrorMsg("This is the last page.");
      setIsLastPage(true);
      return;
    }
    setCursorStack((prev) => [...prev, currentCursor || ""]);
    setCurrentCursor(nextCursor);
    fetchUsers(nextCursor);
    setErrorMsg("");
  }

  function handlePrev() {
    if (cursorStack.length === 0) return;
    const prevStack = [...cursorStack];
    const prevCursor = prevStack.pop() || null;
    setCursorStack(prevStack);
    setCurrentCursor(prevCursor);
    fetchUsers(prevCursor || undefined);
    setErrorMsg("");
    setIsLastPage(false);
  }

  async function updateUser(
    id: string,
    field: "email" | "name",
    value: string
  ) {
    if (!token) return;
    try {
      const res = await fetch("http://localhost:8080/admin/users", {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id,
          [field]: value,
        }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || "Update failed");
      }
      fetchUsers(currentCursor || undefined);
      setErrorMsg("");
    } catch (err: any) {
      setErrorMsg(err.message || "Update failed");
    }
  }

  async function deleteUser(id: string) {
    if (!token) return;
    if (!confirm("Are you sure you want to delete this user?")) return;

    try {
      const res = await fetch("http://localhost:8080/admin/users", {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ id }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || "Delete failed");
      }
      // Refresh list after delete
      fetchUsers(currentCursor || undefined);
      setErrorMsg("");
    } catch (err: any) {
      setErrorMsg(err.message || "Delete failed");
    }
  }

  const handleInlineChange = (
    id: string,
    field: "email" | "name",
    newValue: string
  ) => {
    setUsers((prev) =>
      prev.map((u) => (u.id === id ? { ...u, [field]: newValue } : u))
    );
  };

  const handleInlineBlur = (
    id: string,
    field: "email" | "name",
    value: string
  ) => {
    updateUser(id, field, value);
  };

  return (
    <main className="min-h-screen bg-gray-100 p-8 flex flex-col items-center">
      <div className="mb-4 flex items-center gap-4">
        <h1 className="text-2xl font-semibold">Users</h1>
        <Link
          href="/admin"
          className="rounded bg-gray-600 px-4 py-2 text-white hover:bg-gray-700"
        >
          Back to Admin Page
        </Link>
      </div>

      {errorMsg && <p className="mb-4 text-red-600">{errorMsg}</p>}

      <div className="overflow-x-auto">
        <table className="mx-auto w-full max-w-4xl border border-gray-300 bg-white text-left shadow">
          <thead className="bg-gray-200">
            <tr>
              <th className="border px-4 py-2">User ID</th>
              <th className="border px-4 py-2">Email</th>
              <th className="border px-4 py-2">Name</th>
              <th className="border px-4 py-2">Created At</th>
              <th className="border px-4 py-2">Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((u) => (
              <tr key={u.id} className="hover:bg-gray-50">
                <td className="border px-4 py-2">{u.id}</td>
                <td className="border px-4 py-2">
                  <input
                    type="text"
                    value={u.email}
                    onChange={(e) =>
                      handleInlineChange(u.id, "email", e.target.value)
                    }
                    onBlur={(e) =>
                      handleInlineBlur(u.id, "email", e.target.value)
                    }
                    className="w-full bg-transparent break-all whitespace-normal outline-none"
                  />
                </td>
                <td className="border px-4 py-2">
                  <input
                    type="text"
                    value={u.name}
                    onChange={(e) =>
                      handleInlineChange(u.id, "name", e.target.value)
                    }
                    onBlur={(e) =>
                      handleInlineBlur(u.id, "name", e.target.value)
                    }
                    className="w-full bg-transparent outline-none"
                  />
                </td>
                <td className="border px-4 py-2">
                  {new Date(u.created_at).toLocaleString()}
                </td>
                <td className="border px-4 py-2 text-center">
                  <button
                    onClick={() => deleteUser(u.id)}
                    className="rounded bg-red-500 px-3 py-1 text-white hover:bg-red-600"
                  >
                    Delete
                  </button>
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
