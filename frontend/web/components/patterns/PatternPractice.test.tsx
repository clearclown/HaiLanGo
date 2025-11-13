import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { PatternPractice } from "./PatternPractice";
import type { Pattern, PatternPractice as PatternPracticeType } from "../../lib/types/pattern";

const mockPattern: Pattern = {
	id: "pattern-1",
	book_id: "book-1",
	type: "greeting",
	pattern: "Hello",
	translation: "こんにちは",
	frequency: 5,
	created_at: "2025-01-01T00:00:00Z",
	updated_at: "2025-01-01T00:00:00Z",
};

const mockPractices: PatternPracticeType[] = [
	{
		id: "practice-1",
		pattern_id: "pattern-1",
		question: "How do you say 'Hello' in Japanese?",
		correct_answer: "こんにちは",
		alternative_answers: ["おはよう", "こんばんは", "さようなら"],
		difficulty: 1,
		created_at: "2025-01-01T00:00:00Z",
	},
	{
		id: "practice-2",
		pattern_id: "pattern-1",
		question: "Choose the correct response to 'Hello'",
		correct_answer: "Hello, how are you?",
		alternative_answers: ["Goodbye", "Thank you", "Sorry"],
		difficulty: 2,
		created_at: "2025-01-01T00:00:00Z",
	},
];

describe("PatternPractice", () => {
	it("renders pattern information", () => {
		render(
			<PatternPractice pattern={mockPattern} practices={mockPractices} />,
		);

		expect(screen.getByText("Hello")).toBeInTheDocument();
		expect(screen.getByText("こんにちは")).toBeInTheDocument();
	});

	it("displays practice question", () => {
		render(
			<PatternPractice pattern={mockPattern} practices={mockPractices} />,
		);

		expect(
			screen.getByText("How do you say 'Hello' in Japanese?"),
		).toBeInTheDocument();
	});

	it("shows all answer options", () => {
		render(
			<PatternPractice pattern={mockPattern} practices={mockPractices} />,
		);

		expect(screen.getByText("こんにちは")).toBeInTheDocument();
		expect(screen.getByText("おはよう")).toBeInTheDocument();
		expect(screen.getByText("こんばんは")).toBeInTheDocument();
		expect(screen.getByText("さようなら")).toBeInTheDocument();
	});

	it("marks correct answer when clicked", async () => {
		const onAnswerSubmit = vi.fn();
		render(
			<PatternPractice
				pattern={mockPattern}
				practices={mockPractices}
				onAnswerSubmit={onAnswerSubmit}
			/>,
		);

		const correctAnswer = screen.getByText("こんにちは");
		fireEvent.click(correctAnswer);

		await waitFor(() => {
			expect(onAnswerSubmit).toHaveBeenCalledWith({
				practice_id: "practice-1",
				answer: "こんにちは",
				is_correct: true,
			});
		});
	});

	it("marks incorrect answer when clicked", async () => {
		const onAnswerSubmit = vi.fn();
		render(
			<PatternPractice
				pattern={mockPattern}
				practices={mockPractices}
				onAnswerSubmit={onAnswerSubmit}
			/>,
		);

		const incorrectAnswer = screen.getByText("おはよう");
		fireEvent.click(incorrectAnswer);

		await waitFor(() => {
			expect(onAnswerSubmit).toHaveBeenCalledWith({
				practice_id: "practice-1",
				answer: "おはよう",
				is_correct: false,
			});
		});
	});

	it("moves to next question after answer", async () => {
		render(
			<PatternPractice pattern={mockPattern} practices={mockPractices} />,
		);

		// Answer first question
		const correctAnswer = screen.getByText("こんにちは");
		fireEvent.click(correctAnswer);

		// Wait for next question
		await waitFor(() => {
			expect(
				screen.getByText("Choose the correct response to 'Hello'"),
			).toBeInTheDocument();
		});
	});

	it("shows completion message when all questions answered", async () => {
		render(
			<PatternPractice pattern={mockPattern} practices={mockPractices} />,
		);

		// Answer all questions
		const answers = screen.getAllByRole("button");
		for (const answer of answers.slice(0, 2)) {
			fireEvent.click(answer);
			await waitFor(() => {});
		}

		await waitFor(() => {
			expect(screen.getByText(/completed/i)).toBeInTheDocument();
		});
	});

	it("displays progress indicator", () => {
		render(
			<PatternPractice pattern={mockPattern} practices={mockPractices} />,
		);

		expect(screen.getByText(/1.*2/)).toBeInTheDocument();
	});

	it("shows difficulty level", () => {
		render(
			<PatternPractice pattern={mockPattern} practices={mockPractices} />,
		);

		expect(screen.getByText(/difficulty/i)).toBeInTheDocument();
	});
});
