"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function JoinPage() {
  const [roomId, setRoomId] = useState("");
  const [name, setName] = useState("");
  const router = useRouter();

  const handleJoin = () => {
    if (roomId && name) {
      // Redirect to the game page with roomId and name as query parameter.
      router.push(`/room/${roomId}?name=${encodeURIComponent(name)}`);
    }
  };

  return (
    <div className="flex h-screen items-center justify-center bg-gray-800 text-white">
      <div className="p-6 bg-gray-700 rounded-lg shadow-lg">
        <h1 className="text-2xl font-bold mb-4">Join a Sequence Game Room</h1>
        <input
          type="text"
          placeholder="Enter Room ID"
          className="w-full p-2 mb-2 bg-gray-600 rounded"
          value={roomId}
          onChange={(e) => setRoomId(e.target.value)}
        />
        <input
          type="text"
          placeholder="Enter Your Name"
          className="w-full p-2 mb-4 bg-gray-600 rounded"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <button
          className="w-full p-2 bg-blue-500 rounded hover:bg-blue-600"
          onClick={handleJoin}
        >
          Join Room
        </button>
      </div>
    </div>
  );
}
