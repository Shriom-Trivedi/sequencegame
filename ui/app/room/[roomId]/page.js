"use client";
import { useEffect, useState } from "react";
import { useParams, useSearchParams } from "next/navigation";

export default function GamePage() {
  const { roomId } = useParams();
  const searchParams = useSearchParams();
  const name = searchParams.get("name");

  // Initialize a 10x10 empty board.
  const emptyBoard = Array.from({ length: 10 }, () => Array(10).fill(""));
  const [board, setBoard] = useState(emptyBoard);
  const [turn, setTurn] = useState(null);
  const [ws, setWs] = useState(null);

  useEffect(() => {
    if (!roomId || !name) return;

    const socket = new WebSocket(`ws://localhost:8080/room/${roomId}`);
    socket.onopen = () => {
      console.log("WebSocket connected");
      // On open, send a join message.
      socket.send(JSON.stringify({ action: "join", playerId: name }));
    };

    socket.onmessage = (event) => {
      console.log("Message received:", event.data);
      try {
        const data = JSON.parse(event.data);
        if (data.board && data.currentTurn !== undefined) {
          setBoard(data.board); // This should be an array like you expect, with "Red", "Blue", etc.
          setTurn(data.currentTurn);
        }
      } catch (err) {
        console.error("Error parsing message:", err);
      }
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    setWs(socket);
    return () => socket.close();
  }, [roomId, name]);

  // Handle cell clicks.
  const handleMove = (x, y) => {
    console.log(`Tile clicked at (${x}, ${y})`);
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ action: "move", playerId: name, x, y }));
    } else {
      console.log("WebSocket not open yet.");
    }
  };

  // Determine background color for a cell based on its value.
  const getCellColor = (cell) => {
    console.log({ cell });
    if (cell === "WILD") return "bg-yellow-500";
    if (cell === "Red") return "bg-red-300";
    if (cell === "Blue") return "bg-blue-300";
    return "bg-gray-700";
  };
  console.log({ board });
  return (
    <div className="flex flex-col items-center bg-gray-800 text-white h-screen">
      <h1 className="text-3xl font-bold mt-6">Room: {roomId}</h1>
      <h2 className="text-xl mt-2">Turn: {turn !== null ? turn : "N/A"}</h2>

      <div className="grid grid-cols-10 gap-1 mt-4">
        {board.map((row, rowIdx) =>
          row.map((cell, colIdx) => (
            <div
              key={`${rowIdx}-${colIdx}`}
              onClick={() => handleMove(rowIdx, colIdx)}
              className={`w-10 h-10 flex items-center justify-center border border-gray-600 cursor-pointer ${getCellColor(
                cell
              )}`}
            >
              {/* Optionally display a small label if needed */}
              {cell && cell !== "WILD" && (
                <span className="text-xs text-white font-bold">{cell}</span>
              )}
              {cell === "WILD" && <span className="text-xs font-bold">W</span>}
            </div>
          ))
        )}
      </div>
    </div>
  );
}
