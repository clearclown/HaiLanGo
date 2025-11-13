"use client";

import type React from "react";
import { useState } from "react";
import type {
	Pattern,
	PatternPractice as PatternPracticeType,
} from "../../lib/types/pattern";

interface PatternPracticeProps {
	pattern: Pattern;
	practices: PatternPracticeType[];
	onAnswerSubmit?: (result: {
		practice_id: string;
		answer: string;
		is_correct: boolean;
	}) => void;
	onComplete?: (score: { correct: number; total: number }) => void;
}

export function PatternPractice({
	pattern,
	practices,
	onAnswerSubmit,
	onComplete,
}: PatternPracticeProps) {
	const [currentIndex, setCurrentIndex] = useState(0);
	const [selectedAnswer, setSelectedAnswer] = useState<string | null>(null);
	const [isAnswered, setIsAnswered] = useState(false);
	const [score, setScore] = useState({ correct: 0, total: 0 });

	const currentPractice = practices[currentIndex];
	const isCompleted = currentIndex >= practices.length;

	// Shuffle answers
	const allAnswers = currentPractice
		? [
				currentPractice.correct_answer,
				...currentPractice.alternative_answers,
			].sort(() => Math.random() - 0.5)
		: [];

	const handleAnswerClick = (answer: string) => {
		if (isAnswered) return;

		setSelectedAnswer(answer);
		setIsAnswered(true);

		const isCorrect = answer === currentPractice.correct_answer;
		const newScore = {
			correct: score.correct + (isCorrect ? 1 : 0),
			total: score.total + 1,
		};
		setScore(newScore);

		// Call callback
		if (onAnswerSubmit) {
			onAnswerSubmit({
				practice_id: currentPractice.id,
				answer,
				is_correct: isCorrect,
			});
		}

		// Move to next question after delay
		setTimeout(() => {
			if (currentIndex < practices.length - 1) {
				setCurrentIndex(currentIndex + 1);
				setSelectedAnswer(null);
				setIsAnswered(false);
			} else {
				// Completed all questions
				if (onComplete) {
					onComplete(newScore);
				}
			}
		}, 1500);
	};

	// Completion screen
	if (isCompleted) {
		const percentage = (score.correct / score.total) * 100;

		return (
			<div className="max-w-2xl mx-auto p-8 text-center">
				<h2 className="text-3xl font-bold mb-4">🎉 Practice Completed!</h2>
				<div className="mb-6">
					<div className="text-6xl font-bold text-blue-600 mb-2">
						{percentage.toFixed(0)}%
					</div>
					<p className="text-gray-600">
						{score.correct} out of {score.total} correct
					</p>
				</div>

				<div className="bg-gray-100 rounded-lg p-6 mb-6">
					<h3 className="font-semibold text-lg mb-2">Pattern Practiced</h3>
					<p className="text-2xl font-bold text-gray-900">{pattern.pattern}</p>
					<p className="text-lg text-gray-600">{pattern.translation}</p>
				</div>

				<button
					type="button"
					onClick={() => {
						setCurrentIndex(0);
						setScore({ correct: 0, total: 0 });
						setSelectedAnswer(null);
						setIsAnswered(false);
					}}
					className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
				>
					Practice Again
				</button>
			</div>
		);
	}

	return (
		<div className="max-w-2xl mx-auto p-8">
			{/* Header */}
			<div className="mb-8">
				<div className="flex items-center justify-between mb-4">
					<span className="text-sm text-gray-500">
						Question {currentIndex + 1} of {practices.length}
					</span>
					<span className="text-sm text-gray-500">
						Difficulty: {"⭐".repeat(currentPractice.difficulty)}
					</span>
				</div>

				{/* Progress bar */}
				<div className="w-full bg-gray-200 rounded-full h-2">
					<div
						className="bg-blue-600 h-2 rounded-full transition-all duration-300"
						style={{
							width: `${((currentIndex + 1) / practices.length) * 100}%`,
						}}
					/>
				</div>
			</div>

			{/* Pattern info */}
			<div className="bg-blue-50 rounded-lg p-6 mb-8">
				<p className="text-sm text-gray-600 mb-2">Pattern:</p>
				<p className="text-2xl font-bold text-gray-900">{pattern.pattern}</p>
				<p className="text-lg text-gray-600">{pattern.translation}</p>
			</div>

			{/* Question */}
			<div className="mb-8">
				<h2 className="text-2xl font-semibold mb-6">
					{currentPractice.question}
				</h2>

				{/* Answer options */}
				<div className="space-y-3">
					{allAnswers.map((answer, index) => {
						const isCorrect = answer === currentPractice.correct_answer;
						const isSelected = answer === selectedAnswer;

						let buttonClass =
							"w-full p-4 rounded-lg border-2 text-left transition-all ";

						if (isAnswered) {
							if (isSelected) {
								buttonClass += isCorrect
									? "bg-green-100 border-green-500 text-green-900"
									: "bg-red-100 border-red-500 text-red-900";
							} else if (isCorrect) {
								buttonClass += "bg-green-100 border-green-500 text-green-900";
							} else {
								buttonClass += "border-gray-200 text-gray-400";
							}
						} else {
							buttonClass +=
								"border-gray-200 hover:border-blue-400 hover:bg-blue-50 cursor-pointer";
						}

						return (
							<button
								key={`${currentPractice.id}-${index.toString()}`}
								type="button"
								onClick={() => handleAnswerClick(answer)}
								disabled={isAnswered}
								className={buttonClass}
							>
								<span className="font-medium">{answer}</span>
								{isAnswered && isSelected && (
									<span className="ml-2">
										{isCorrect ? "✓" : "✗"}
									</span>
								)}
							</button>
						);
					})}
				</div>
			</div>

			{/* Feedback */}
			{isAnswered && (
				<div
					className={`p-4 rounded-lg ${
						selectedAnswer === currentPractice.correct_answer
							? "bg-green-100 text-green-900"
							: "bg-red-100 text-red-900"
					}`}
				>
					{selectedAnswer === currentPractice.correct_answer ? (
						<p className="font-semibold">🎉 Correct! Well done!</p>
					) : (
						<p className="font-semibold">
							Not quite. The correct answer is:{" "}
							{currentPractice.correct_answer}
						</p>
					)}
				</div>
			)}
		</div>
	);
}
