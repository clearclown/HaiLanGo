"use client";

import type React from "react";
import type { Pattern, PatternType } from "../../lib/types/pattern";

interface PatternListProps {
	patterns: Pattern[];
	onPatternClick?: (pattern: Pattern) => void;
	onFilterChange?: (type: PatternType | "all") => void;
	sortBy?: "frequency" | "type" | "recent";
}

const PATTERN_TYPE_LABELS: Record<PatternType, string> = {
	greeting: "Greeting",
	question: "Question",
	response: "Response",
	request: "Request",
	confirmation: "Confirmation",
	other: "Other",
};

const PATTERN_TYPE_COLORS: Record<PatternType, string> = {
	greeting: "bg-blue-100 text-blue-800",
	question: "bg-purple-100 text-purple-800",
	response: "bg-green-100 text-green-800",
	request: "bg-yellow-100 text-yellow-800",
	confirmation: "bg-indigo-100 text-indigo-800",
	other: "bg-gray-100 text-gray-800",
};

export function PatternList({
	patterns,
	onPatternClick,
	onFilterChange,
	sortBy = "frequency",
}: PatternListProps) {
	// Sort patterns
	const sortedPatterns = [...patterns].sort((a, b) => {
		if (sortBy === "frequency") {
			return b.frequency - a.frequency;
		}
		if (sortBy === "type") {
			return a.type.localeCompare(b.type);
		}
		return (
			new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
		);
	});

	// Empty state
	if (patterns.length === 0) {
		return (
			<div className="text-center py-12">
				<p className="text-gray-500 text-lg">No patterns found</p>
				<p className="text-gray-400 text-sm mt-2">
					Try extracting patterns from a book first
				</p>
			</div>
		);
	}

	return (
		<div className="space-y-4">
			{/* Filter buttons */}
			{onFilterChange && (
				<div className="flex flex-wrap gap-2 mb-6">
					<button
						type="button"
						onClick={() => onFilterChange("all")}
						className="px-4 py-2 rounded-lg bg-gray-100 hover:bg-gray-200 text-gray-700 text-sm font-medium transition-colors"
					>
						All
					</button>
					{Object.entries(PATTERN_TYPE_LABELS).map(([type, label]) => (
						<button
							key={type}
							type="button"
							onClick={() => onFilterChange(type as PatternType)}
							className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
								PATTERN_TYPE_COLORS[type as PatternType]
							}`}
						>
							{label}
						</button>
					))}
				</div>
			)}

			{/* Pattern cards */}
			<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{sortedPatterns.map((pattern) => (
					<div
						key={pattern.id}
						data-testid="pattern-card"
						onClick={() => onPatternClick?.(pattern)}
						onKeyDown={(e) => {
							if (e.key === "Enter" || e.key === " ") {
								onPatternClick?.(pattern);
							}
						}}
						className={`
              p-6 rounded-lg border-2 border-gray-200
              hover:border-blue-400 hover:shadow-lg
              transition-all cursor-pointer
              ${onPatternClick ? "hover:scale-105" : ""}
            `}
						role="button"
						tabIndex={0}
					>
						<div className="flex items-start justify-between mb-3">
							<span
								className={`
                  px-3 py-1 rounded-full text-xs font-semibold
                  ${PATTERN_TYPE_COLORS[pattern.type]}
                `}
							>
								{PATTERN_TYPE_LABELS[pattern.type]}
							</span>
							<span className="text-gray-500 text-sm">
								×{pattern.frequency}
							</span>
						</div>

						<div className="space-y-2">
							<p className="text-lg font-semibold text-gray-900">
								{pattern.pattern}
							</p>
							<p className="text-base text-gray-600">{pattern.translation}</p>
						</div>
					</div>
				))}
			</div>
		</div>
	);
}
