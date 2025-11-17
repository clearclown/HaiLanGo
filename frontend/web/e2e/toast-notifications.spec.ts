import { test, expect, type Page } from "@playwright/test";

/**
 * Toast Notifications E2E Tests
 *
 * This test suite verifies the toast notification system integrated with WebSocket.
 * Tests use window.__TEST_TOAST__ API exposed by ToastProvider for E2E testing.
 */

test.describe("Toast Notification System", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.waitForLoadState("networkidle");

    // Wait for toast system to initialize
    await page.waitForFunction(() => {
      return typeof (window as any).__TEST_TOAST__ !== 'undefined';
    }, { timeout: 5000 });
  });

  test("should render toast notification when triggered", async ({ page }) => {
    // Trigger a notification using the test API
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showInfo("Test Notification", "This is a test message");
    });

    // Wait for toast to appear using data-testid
    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible({ timeout: 5000 });

    // Verify toast content
    await expect(toast).toContainText("Test Notification");
    await expect(toast).toContainText("This is a test message");
  });

  test("should display different toast types with correct styling", async ({ page }) => {
    const types = [
      { type: "info", icon: "🔵" },
      { type: "success", icon: "✅" },
      { type: "warning", icon: "⚠️" },
      { type: "error", icon: "❌" },
    ];

    for (const { type, icon } of types) {
      await page.evaluate((notifType) => {
        const toast = (window as any).__TEST_TOAST__;
        const methodName = `show${notifType.charAt(0).toUpperCase()}${notifType.slice(1)}` as keyof typeof toast;
        toast[methodName](`${notifType.toUpperCase()} Test`, `Testing ${notifType} notification`);
      }, type);

      // Wait for toast to appear
      const toast = page.getByTestId("toast-notification").last();
      await expect(toast).toBeVisible({ timeout: 5000 });

      // Verify icon and type attribute
      await expect(toast).toContainText(icon);
      await expect(toast).toHaveAttribute("data-toast-type", type);

      // Verify title
      await expect(toast).toContainText(`${type.toUpperCase()} Test`);
    }
  });

  test("should auto-dismiss toast after default duration (5 seconds)", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showInfo("Auto-dismiss Test", "This should disappear in 5 seconds");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible({ timeout: 5000 });

    // Wait for auto-dismiss (5 seconds + 300ms animation)
    await page.waitForTimeout(5500);

    // Toast should be hidden or removed
    await expect(toast).not.toBeVisible();
  });

  test("should handle multiple toast notifications", async ({ page }) => {
    // Trigger 3 notifications
    await page.evaluate(() => {
      const toast = (window as any).__TEST_TOAST__;
      toast.showInfo("Notification 1", "Message 1");
      toast.showInfo("Notification 2", "Message 2");
      toast.showInfo("Notification 3", "Message 3");
    });

    // All 3 toasts should be visible
    const toasts = page.getByTestId("toast-notification");
    await expect(toasts).toHaveCount(3, { timeout: 5000 });

    // Verify content of each
    await expect(toasts.nth(0)).toContainText("Notification 1");
    await expect(toasts.nth(1)).toContainText("Notification 2");
    await expect(toasts.nth(2)).toContainText("Notification 3");
  });

  test("should close toast when close button is clicked", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showInfo("Closeable Toast", "Click the X to close");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible({ timeout: 5000 });

    // Click close button (SVG button)
    const closeButton = toast.locator('button[aria-label="Close"]');
    await closeButton.click();

    // Toast should disappear
    await expect(toast).not.toBeVisible({ timeout: 1000 });
  });
});

test.describe("Toast Notification Edge Cases", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.waitForLoadState("networkidle");
    await page.waitForFunction(() => typeof (window as any).__TEST_TOAST__ !== 'undefined');
  });

  test("should handle rapid consecutive notifications", async ({ page }) => {
    // Trigger 10 notifications rapidly
    await page.evaluate(() => {
      const toast = (window as any).__TEST_TOAST__;
      for (let i = 1; i <= 10; i++) {
        toast.showInfo(`Rapid Notification ${i}`, `Message ${i}`);
      }
    });

    // Should display all notifications
    const toasts = page.getByTestId("toast-notification");
    const count = await toasts.count();

    // All 10 should be visible
    expect(count).toBe(10);
  });

  test("should handle very long notification messages", async ({ page }) => {
    const longMessage = "A".repeat(500);

    await page.evaluate((msg: string) => {
      (window as any).__TEST_TOAST__.showInfo("Long Message Test", msg);
    }, longMessage);

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible({ timeout: 5000 });

    // Toast should be visible and not break the layout
    const boundingBox = await toast.boundingBox();
    expect(boundingBox).not.toBeNull();
    if (boundingBox) {
      expect(boundingBox.width).toBeLessThanOrEqual(500); // max-w-[500px]
    }
  });

  test("should handle notifications with missing optional fields", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showInfo("Minimal Notification");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible({ timeout: 5000 });

    // Should still display with just the title
    await expect(toast).toContainText("Minimal Notification");
  });

  test("should position toasts in top-right corner", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showInfo("Position Test", "Checking position");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible({ timeout: 5000 });

    const boundingBox = await toast.boundingBox();
    expect(boundingBox).not.toBeNull();

    // Should be in top-right area
    if (boundingBox) {
      const viewportSize = page.viewportSize();
      if (viewportSize) {
        // Right side of screen
        expect(boundingBox.x).toBeGreaterThan(viewportSize.width / 2);
        // Top of screen
        expect(boundingBox.y).toBeLessThan(200);
      }
    }
  });
});

test.describe("Toast Accessibility", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.waitForLoadState("networkidle");
    await page.waitForFunction(() => typeof (window as any).__TEST_TOAST__ !== 'undefined');
  });

  test("should have proper ARIA role for accessibility", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showInfo("Accessibility Test", "Testing ARIA attributes");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible({ timeout: 5000 });

    // Verify ARIA attributes
    await expect(toast).toHaveAttribute("role", "alert");
    await expect(toast).toHaveAttribute("aria-live", "assertive");
  });

  test("should be keyboard accessible (close button focusable)", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showInfo("Keyboard Test", "Testing keyboard navigation");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible({ timeout: 5000 });

    // Focus on close button
    const closeButton = toast.locator('button[aria-label="Close"]');
    await closeButton.focus();

    // Verify button is focused
    await expect(closeButton).toBeFocused();

    // Press Enter to close
    await page.keyboard.press("Enter");

    // Toast should close
    await expect(toast).not.toBeVisible({ timeout: 1000 });
  });
});

test.describe("Integration with Toast Types", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.waitForLoadState("networkidle");
    await page.waitForFunction(() => typeof (window as any).__TEST_TOAST__ !== 'undefined');
  });

  test("should use showInfo correctly", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showInfo("Info Title", "Info message");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible();
    await expect(toast).toHaveAttribute("data-toast-type", "info");
    await expect(toast).toContainText("🔵");
  });

  test("should use showSuccess correctly", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showSuccess("Success Title", "Success message");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible();
    await expect(toast).toHaveAttribute("data-toast-type", "success");
    await expect(toast).toContainText("✅");
  });

  test("should use showWarning correctly", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showWarning("Warning Title", "Warning message");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible();
    await expect(toast).toHaveAttribute("data-toast-type", "warning");
    await expect(toast).toContainText("⚠️");
  });

  test("should use showError correctly", async ({ page }) => {
    await page.evaluate(() => {
      (window as any).__TEST_TOAST__.showError("Error Title", "Error message");
    });

    const toast = page.getByTestId("toast-notification");
    await expect(toast).toBeVisible();
    await expect(toast).toHaveAttribute("data-toast-type", "error");
    await expect(toast).toContainText("❌");
  });
});
