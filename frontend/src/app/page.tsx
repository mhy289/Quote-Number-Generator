"use client";

import { useState } from "react";

export default function Home() {
    // 随机数
    const [min, setMin] = useState(1);
    const [max, setMax] = useState(100);
    const [randomNumber, setRandomNumber] = useState<number | null>(null);
    const [numberLoading, setNumberLoading] = useState(false);
    const [numberError, setNumberError] = useState("");

    // 名言
    const [quote, setQuote] = useState("");
    const [quoteLoading, setQuoteLoading] = useState(false);
    const [quoteError, setQuoteError] = useState("");

    // 调用后端生成随机数
    const handleGenerateNumber = async () => {
        setNumberLoading(true);
        setNumberError("");
        setRandomNumber(null);
        try {
            const res = await fetch("http://localhost:8080/generator.GeneratorService/GetRandomNumber", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ min, max }),
            });
            if (!res.ok) {
                const err = await res.json();
                setNumberError(err.message || "请求失败");
                return;
            }
            const data = await res.json();
            setRandomNumber(data.number);
        } catch {
            setNumberError("网络异常，请检查后端服务是否启动");
        } finally {
            setNumberLoading(false);
        }
    };

    // 调用后端获取随机名言
    const handleGetQuote = async () => {
        setQuoteLoading(true);
        setQuoteError("");
        setQuote("");
        try {
            const res = await fetch("http://localhost:8080/generator.GeneratorService/GetRandomQuote", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({}),
            });
            if (!res.ok) {
                const err = await res.json();
                setQuoteError(err.message || "请求失败");
                return;
            }
            const data = await res.json();
            setQuote(data.quote);
        } catch {
            setQuoteError("网络异常，请检查后端服务是否启动");
        } finally {
            setQuoteLoading(false);
        }
    };

    return (
        <div className="p-8 space-y-8 max-w-lg mx-auto">
            <h1 className="text-2xl font-bold text-center">QuoteGenerator</h1>

            <section className="space-y-4 p-6 border rounded-lg shadow-sm">
                <h2 className="text-xl font-bold">随机数生成器</h2>
                <div className="flex flex-wrap gap-2">
                    <input
                        type="number"
                        value={min}
                        onChange={(e) => setMin(Number(e.target.value))}
                        placeholder="Min"
                        className="border p-2 rounded flex-1 min-w-[80px]"
                    />
                    <input
                        type="number"
                        value={max}
                        onChange={(e) => setMax(Number(e.target.value))}
                        placeholder="Max"
                        className="border p-2 rounded flex-1 min-w-[80px]"
                    />
                    <button
                        onClick={handleGenerateNumber}
                        disabled={numberLoading}
                        className="px-4 py-2 bg-blue-500 text-white rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-blue-600 transition-colors"
                    >
                        {numberLoading ? "生成中..." : "生成"}
                    </button>
                </div>
                {numberError && <p className="text-red-500 text-sm">{numberError}</p>}
                {randomNumber !== null && (
                    <p className="text-lg font-semibold">随机数：<span className="text-blue-600">{randomNumber}</span></p>
                )}
            </section>

            <section className="space-y-4 p-6 border rounded-lg shadow-sm">
                <h2 className="text-xl font-bold">随机名言</h2>
                <button
                    onClick={handleGetQuote}
                    disabled={quoteLoading}
                    className="px-4 py-2 bg-green-500 text-white rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-green-600 transition-colors"
                >
                    {quoteLoading ? "获取中..." : "获取名言"}
                </button>
                {quoteError && <p className="text-red-500 text-sm">{quoteError}</p>}
                {quote && (
                    <blockquote className="italic border-l-4 border-green-500 pl-4 py-2 text-gray-700">
                        &ldquo;{quote}&rdquo;
                    </blockquote>
                )}
            </section>
        </div>
    );
}