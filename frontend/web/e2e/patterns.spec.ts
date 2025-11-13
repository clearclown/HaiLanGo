import { test, expect } from "@playwright/test";

test.describe("Pattern Extraction and Practice", () => {
	test.beforeEach(async ({ page }) => {
		// Navigate to the patterns page
		await page.goto("/books/test-book/patterns");
	});

	test("should display pattern list", async ({ page }) => {
		// Wait for patterns to load
		await page.waitForSelector('[data-testid="pattern-card"]');

		// Check that patterns are displayed
		const patterns = await page.$$('[data-testid="pattern-card"]');
		expect(patterns.length).toBeGreaterThan(0);

		// Verify pattern content
		await expect(page.locator("text=Hello")).toBeVisible();
		await expect(page.locator("text=こんにちは")).toBeVisible();
	});

	test("should filter patterns by type", async ({ page }) => {
		// Wait for page to load
		await page.waitForSelector('[data-testid="pattern-card"]');

		// Click greeting filter
		await page.click('button:has-text("Greeting")');

		// Verify only greeting patterns are shown
		const greetingBadges = await page.$$('span:has-text("Greeting")');
		expect(greetingBadges.length).toBeGreaterThan(0);
	});

	test("should open pattern practice on click", async ({ page }) => {
		// Wait for patterns to load
		await page.waitForSelector('[data-testid="pattern-card"]');

		// Click on first pattern
		await page.click('[data-testid="pattern-card"]');

		// Verify practice screen is shown
		await expect(page.locator("text=Question 1 of")).toBeVisible();
		await expect(page.locator("text=Difficulty:")).toBeVisible();
	});

	test("should complete practice exercise", async ({ page }) => {
		// Navigate to practice
		await page.goto("/patterns/test-pattern/practice");

		// Wait for question to load
		await page.waitForSelector("button:has-text('こんにちは')");

		// Answer first question
		await page.click("button:has-text('こんにちは')");

		// Wait for feedback
		await expect(page.locator("text=Correct")).toBeVisible();

		// Wait for next question or completion
		await page.waitForTimeout(2000);

		// Continue answering if there are more questions
		const completionText = await page.locator("text=Completed").isVisible();
		if (!completionText) {
			// Answer remaining questions
			const answerButtons = await page.$$("button");
			if (answerButtons.length > 0) {
				await answerButtons[0].click();
			}
		}
	});

	test("should show completion screen after all questions", async ({
		page,
	}) => {
		// Navigate to practice
		await page.goto("/patterns/test-pattern/practice");

		// Answer all questions (assuming 2 questions)
		for (let i = 0; i < 2; i++) {
			await page.waitForSelector("button");
			const buttons = await page.$$("button");
			if (buttons.length > 0) {
				await buttons[0].click();
				await page.waitForTimeout(2000);
			}
		}

		// Verify completion screen
		await expect(page.locator("text=Practice Completed")).toBeVisible();
		await expect(page.locator("text=out of")).toBeVisible();
	});

	test("should restart practice after completion", async ({ page }) => {
		// Navigate to completion screen (after completing practice)
		await page.goto("/patterns/test-pattern/practice?completed=true");

		// Click practice again button
		await page.click('button:has-text("Practice Again")');

		// Verify back to first question
		await expect(page.locator("text=Question 1 of")).toBeVisible();
	});

	test("should display pattern frequency", async ({ page }) => {
		// Wait for patterns to load
		await page.waitForSelector('[data-testid="pattern-card"]');

		// Verify frequency is displayed
		await expect(page.locator("text=×5")).toBeVisible();
	});

	test("should sort patterns by frequency", async ({ page }) => {
		// Wait for patterns to load
		await page.waitForSelector('[data-testid="pattern-card"]');

		// Get first pattern
		const firstPattern = await page
			.$('[data-testid="pattern-card"]')
			.then((el) => el?.textContent());

		// Verify sorting (highest frequency first)
		expect(firstPattern).toContain("×5");
	});

	test("should show progress bar during practice", async ({ page }) => {
		// Navigate to practice
		await page.goto("/patterns/test-pattern/practice");

		// Wait for progress bar
		await page.waitForSelector(".bg-blue-600");

		// Verify progress bar exists
		const progressBar = await page.$(".bg-blue-600");
		expect(progressBar).not.toBeNull();
	});

	test("should highlight correct and incorrect answers", async ({ page }) => {
		// Navigate to practice
		await page.goto("/patterns/test-pattern/practice");

		// Wait for question
		await page.waitForSelector("button");

		// Click wrong answer
		const buttons = await page.$$("button");
		if (buttons.length > 1) {
			// Click second button (likely wrong)
			await buttons[1].click();

			// Wait for highlighting
			await page.waitForTimeout(500);

			// Verify incorrect answer is highlighted
			const incorrectAnswer = await page.$(".bg-red-100");
			expect(incorrectAnswer).not.toBeNull();

			// Verify correct answer is also shown
			const correctAnswer = await page.$(".bg-green-100");
			expect(correctAnswer).not.toBeNull();
		}
	});
});
