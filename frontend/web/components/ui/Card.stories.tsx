import type { Meta, StoryObj } from '@storybook/react';
import { Badge } from './Badge';
import { Button } from './Button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardImage,
  CardTitle,
} from './Card';
import { Progress } from './Progress';

const meta = {
  title: 'UI/Card',
  component: Card,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component:
          'Card component following HaiLanGo design system. Border-radius: 12px, shadow: 0 2px 8px rgba(0,0,0,0.06).',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'bordered', 'elevated', 'interactive', 'highlight'],
      description: 'Card style variant',
    },
    padding: {
      control: 'select',
      options: ['none', 'sm', 'md', 'lg'],
      description: 'Card padding size',
    },
  },
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

// Basic Card
export const Default: Story = {
  args: {
    variant: 'default',
    padding: 'md',
    children: (
      <>
        <CardHeader>
          <CardTitle>Card Title</CardTitle>
          <CardDescription>This is a card description with some details.</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-text-primary">Card content goes here.</p>
        </CardContent>
      </>
    ),
  },
};

// Bordered Card
export const Bordered: Story = {
  args: {
    variant: 'bordered',
    padding: 'md',
    children: (
      <>
        <CardHeader>
          <CardTitle>Bordered Card</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-text-secondary">A card with a subtle border instead of shadow.</p>
        </CardContent>
      </>
    ),
  },
};

// Elevated Card
export const Elevated: Story = {
  args: {
    variant: 'elevated',
    padding: 'md',
    children: (
      <>
        <CardHeader>
          <CardTitle>Elevated Card</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-text-secondary">A card with stronger shadow for emphasis.</p>
        </CardContent>
      </>
    ),
  },
};

// Interactive Card
export const Interactive: Story = {
  args: {
    variant: 'interactive',
    padding: 'md',
    children: (
      <>
        <CardHeader>
          <CardTitle>Interactive Card</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-text-secondary">Hover to see the interaction effect.</p>
        </CardContent>
      </>
    ),
  },
};

// Highlight Card
export const Highlight: Story = {
  args: {
    variant: 'highlight',
    padding: 'md',
    children: (
      <>
        <CardHeader>
          <CardTitle>Highlighted Card</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-text-secondary">
            A card with gradient background for special content.
          </p>
        </CardContent>
      </>
    ),
  },
};

// Card with Footer
export const WithFooter: Story = {
  args: {
    variant: 'default',
    padding: 'md',
    children: (
      <>
        <CardHeader>
          <CardTitle>Card with Actions</CardTitle>
          <CardDescription>A card with action buttons in the footer.</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-text-primary">Main content area of the card.</p>
        </CardContent>
        <CardFooter variant="actions">
          <Button variant="primary" size="sm">
            Save
          </Button>
          <Button variant="ghost" size="sm">
            Cancel
          </Button>
        </CardFooter>
      </>
    ),
  },
};

// Card with Header Action
export const WithHeaderAction: Story = {
  args: {
    variant: 'default',
    padding: 'md',
    children: (
      <>
        <CardHeader
          action={
            <Button variant="ghost" size="sm">
              Edit
            </Button>
          }
        >
          <CardTitle>Editable Card</CardTitle>
          <CardDescription>Card with action button in header.</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-text-primary">Content that can be edited.</p>
        </CardContent>
      </>
    ),
  },
};

// Book Card (HaiLanGo specific)
export const BookCard: Story = {
  render: () => (
    <Card variant="interactive" padding="md" className="w-72">
      <div className="flex gap-4">
        <div className="w-16 h-20 bg-gradient-to-br from-primary to-secondary rounded-lg flex items-center justify-center flex-shrink-0">
          <span className="text-white text-2xl">📕</span>
        </div>
        <div className="flex-1 min-w-0">
          <CardTitle className="truncate">Russian Language 101</CardTitle>
          <CardDescription className="truncate">Introduction to Russian</CardDescription>
          <div className="mt-2">
            <Progress value={45} size="sm" showLabel />
          </div>
        </div>
      </div>
      <CardFooter variant="bordered" className="flex justify-between items-center">
        <span className="text-xs text-text-secondary">Last studied: 2 hours ago</span>
        <Badge variant="primary" size="sm">
          45%
        </Badge>
      </CardFooter>
    </Card>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Today's Learning Card
export const TodayLearningCard: Story = {
  render: () => (
    <Card variant="highlight" padding="lg" className="w-80">
      <CardHeader>
        <div className="flex items-center gap-2">
          <span className="text-2xl">📚</span>
          <CardTitle>Today's Learning</CardTitle>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div>
            <div className="flex justify-between text-sm mb-1">
              <span className="text-text-secondary">Progress</span>
              <span className="font-medium text-text-primary">2/5 pages</span>
            </div>
            <Progress value={40} variant="gradient" />
          </div>
          <Button variant="primary" size="lg" fullWidth>
            Continue Learning
          </Button>
        </div>
      </CardContent>
    </Card>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Stats Card
export const StatsCard: Story = {
  render: () => (
    <div className="grid grid-cols-2 gap-4">
      <Card variant="default" padding="md" className="text-center">
        <div className="text-3xl mb-1">📖</div>
        <div className="text-2xl font-bold text-text-primary">5</div>
        <div className="text-sm text-text-secondary">Books</div>
      </Card>
      <Card variant="default" padding="md" className="text-center">
        <div className="text-3xl mb-1">🎯</div>
        <div className="text-2xl font-bold text-text-primary">12</div>
        <div className="text-sm text-text-secondary">Reviews Due</div>
      </Card>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Streak Card
export const StreakCard: Story = {
  render: () => (
    <Card variant="elevated" padding="md" className="w-72">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Learning Streak</CardTitle>
          <span className="text-2xl">🔥</span>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="flex items-baseline gap-2">
            <span className="text-4xl font-bold text-primary">7</span>
            <span className="text-text-secondary">days</span>
          </div>
          <p className="text-sm text-text-secondary">Keep going! Your longest streak is 15 days.</p>
        </div>
      </CardContent>
    </Card>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Phrase Card (for learning)
export const PhraseCard: Story = {
  render: () => (
    <Card variant="bordered" padding="lg" className="w-96">
      <div className="text-center space-y-4">
        <div className="text-3xl font-medium text-text-primary">Здравствуйте!</div>
        <div className="text-lg text-text-secondary">Hello! (formal)</div>
        <div className="pt-4 flex justify-center gap-3">
          <Button
            variant="outline"
            size="md"
            leftIcon={
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15.536a5 5 0 001.414 1.414m2.828-9.9a9 9 0 012.828-2.828"
                />
              </svg>
            }
          >
            Listen
          </Button>
          <Button
            variant="primary"
            size="md"
            leftIcon={
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"
                />
              </svg>
            }
          >
            Speak
          </Button>
        </div>
      </div>
    </Card>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Pronunciation Score Card
export const PronunciationScoreCard: Story = {
  render: () => (
    <Card variant="default" padding="lg" className="w-80 text-center">
      <div className="space-y-4">
        <div className="inline-flex items-center justify-center w-24 h-24 rounded-full bg-success/10 text-success">
          <span className="text-4xl font-bold">85</span>
        </div>
        <div>
          <CardTitle>Great pronunciation!</CardTitle>
          <CardDescription>Just a bit more practice on the ending sound.</CardDescription>
        </div>
        <div className="space-y-2 text-left">
          <div className="flex justify-between text-sm">
            <span className="text-text-secondary">Accuracy</span>
            <span className="text-text-primary font-medium">85%</span>
          </div>
          <Progress value={85} variant="success" size="sm" />
          <div className="flex justify-between text-sm">
            <span className="text-text-secondary">Fluency</span>
            <span className="text-text-primary font-medium">78%</span>
          </div>
          <Progress value={78} variant="primary" size="sm" />
        </div>
        <div className="flex gap-2 pt-2">
          <Button variant="outline" size="md" className="flex-1">
            Try Again
          </Button>
          <Button variant="primary" size="md" className="flex-1">
            Next
          </Button>
        </div>
      </div>
    </Card>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// All Variants
export const AllVariants: Story = {
  render: () => (
    <div className="grid grid-cols-2 gap-4">
      <Card variant="default" padding="md">
        <CardTitle>Default</CardTitle>
        <CardDescription>Subtle shadow</CardDescription>
      </Card>
      <Card variant="bordered" padding="md">
        <CardTitle>Bordered</CardTitle>
        <CardDescription>Border style</CardDescription>
      </Card>
      <Card variant="elevated" padding="md">
        <CardTitle>Elevated</CardTitle>
        <CardDescription>Strong shadow</CardDescription>
      </Card>
      <Card variant="interactive" padding="md">
        <CardTitle>Interactive</CardTitle>
        <CardDescription>Hover effect</CardDescription>
      </Card>
      <Card variant="highlight" padding="md" className="col-span-2">
        <CardTitle>Highlight</CardTitle>
        <CardDescription>Gradient background for emphasis</CardDescription>
      </Card>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};
