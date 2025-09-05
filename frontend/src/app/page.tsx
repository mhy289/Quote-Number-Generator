"use client";

import { useState } from "react";

export default function Home() {
  // 随机数
  const [min, setMin] = useState(1);
  const [max, setMax] = useState(100);
  const [randomNumber, setRandomNumber] = useState<number | null>(null);

  // 名言
  const [quote, setQuote] = useState("");

  // 调用后端生成随机数
  const handleGenerateNumber = async () => {
    const res = await fetch("http://localhost:8080/generator.GeneratorService/GetRandomNumber", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ min, max }),
    });
    const data = await res.json();
    setRandomNumber(data.number);
  };

  // 调用后端获取随机名言
  const handleGetQuote = async () => {
    const res = await fetch("http://localhost:8080/generator.GeneratorService/GetRandomQuote", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({}),
    });
    const data = await res.json();
    setQuote(data.quote);
  };

  return (
    <div className="p-8 space-y-8">
      <section className="space-y-4">
        <h2 className="text-xl font-bold">随机数生成器</h2>
        <input
          type="number"
          value={min}
          onChange={(e) => setMin(Number(e.target.value))}
          placeholder="Min"
          className="border p-2 rounded mr-2"
        />
        <input
          type="number"
          value={max}
          onChange={(e) => setMax(Number(e.target.value))}
          placeholder="Max"
          className="border p-2 rounded mr-2"
        />
        <button onClick={handleGenerateNumber} className="px-4 py-2 bg-blue-500 text-white rounded">
          生成
        </button>
        {randomNumber !== null && <p>随机数：{randomNumber}</p>}
      </section>

      <section className="space-y-4">
        <h2 className="text-xl font-bold">随机名言</h2>
        <button onClick={handleGetQuote} className="px-4 py-2 bg-green-500 text-white rounded">
          获取名言
        </button>
        {quote && <p>名言：{quote}</p>}
      </section>
    </div>
  );
}
