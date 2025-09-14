"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

export default function Home() {
  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  // ✅ Auto-redirect if user already signed in
  useEffect(() => {
    const storedEmail = localStorage.getItem("email");
    const storedToken = localStorage.getItem("token");
    if (storedEmail && storedToken) {
      router.replace("/admin"); // replace avoids back button returning to login
    }
  }, [router]);

  async function handleSignup() {
    if (!email || !password) {
      setMessage("Enter email & password");
      return;
    }
    setLoading(true);
    setMessage("");

    try {
      const res = await fetch("http://localhost:8080/users/signup", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) throw new Error(`Signup failed (${res.status})`);

      const data = await res.json();
      setMessage(`Account created! User ID: ${data.user_id}`);
    } catch (err: any) {
      setMessage(err.message || "Signup failed");
    } finally {
      setLoading(false);
    }
  }

  async function handleSignin() {
    if (!email || !password) {
      setMessage("Enter email & password");
      return;
    }
    setLoading(true);
    setMessage("");

    try {
      const res = await fetch("http://localhost:8080/users/signin", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) throw new Error(`Signin failed (${res.status})`);

      const data = await res.json();
      localStorage.setItem("email", email);
      localStorage.setItem("token", data.token);
      router.push("/admin");
    } catch (err: any) {
      setMessage(err.message || "Signin failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="flex h-screen items-center justify-center bg-gray-100">
      <div className="w-full max-w-xs rounded-2xl bg-white p-6 shadow">
        <h1 className="mb-4 text-xl font-semibold text-center">Login</h1>

        <input
          type="email"
          placeholder="Email"
          className="mb-2 w-full rounded border p-2"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          className="mb-4 w-full rounded border p-2"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        <div className="flex gap-2">
          <button
            onClick={handleSignin}
            disabled={loading}
            className="flex-1 rounded bg-blue-500 p-2 text-white hover:bg-blue-600 disabled:opacity-50"
          >
            {loading ? "Loading..." : "Login"}
          </button>
          <button
            onClick={handleSignup}
            disabled={loading}
            className="flex-1 rounded bg-green-500 p-2 text-white hover:bg-green-600 disabled:opacity-50"
          >
            {loading ? "Loading..." : "Sign Up"}
          </button>
        </div>

        {message && (
          <p className="mt-4 text-center text-sm text-gray-700">{message}</p>
        )}
      </div>
    </main>
  );
}
