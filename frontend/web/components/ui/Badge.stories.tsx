import type { Meta, StoryObj } from '@storybook/react';
import { Badge, BadgeGroup } from './Badge';

const meta = {
  title: 'UI/Badge',
  component: Badge,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component:
          'Badge component for status indicators, labels, and tags. Used for learning progress, language labels, etc.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'primary', 'secondary', 'success', 'warning', 'danger', 'info'],
      description: 'Badge color variant',
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
      description: 'Badge size',
    },
    dot: {
      control: 'boolean',
      description: 'Show status dot',
    },
  },
} satisfies Meta<typeof Badge>;

export default meta;
type Story = StoryObj<typeof meta>;

// Default Badge
export const Default: Story = {
  args: {
    children: 'Badge',
    variant: 'default',
    size: 'sm',
  },
};

// Primary Badge
export const Primary: Story = {
  args: {
    children: 'Primary',
    variant: 'primary',
  },
};

// Secondary Badge
export const Secondary: Story = {
  args: {
    children: 'Secondary',
    variant: 'secondary',
  },
};

// Success Badge
export const Success: Story = {
  args: {
    children: 'Success',
    variant: 'success',
  },
};

// Warning Badge
export const Warning: Story = {
  args: {
    children: 'Warning',
    variant: 'warning',
  },
};

// Danger Badge
export const Danger: Story = {
  args: {
    children: 'Danger',
    variant: 'danger',
  },
};

// Info Badge
export const Info: Story = {
  args: {
    children: 'Info',
    variant: 'info',
  },
};

// With Dot
export const WithDot: Story = {
  args: {
    children: 'Online',
    variant: 'success',
    dot: true,
  },
};

// With Icon
const CheckIcon = () => (
  <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
  </svg>
);

export const WithIcon: Story = {
  args: {
    children: 'Completed',
    variant: 'success',
    icon: <CheckIcon />,
  },
};

// All Variants
export const AllVariants: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <Badge variant="default">Default</Badge>
      <Badge variant="primary">Primary</Badge>
      <Badge variant="secondary">Secondary</Badge>
      <Badge variant="success">Success</Badge>
      <Badge variant="warning">Warning</Badge>
      <Badge variant="danger">Danger</Badge>
      <Badge variant="info">Info</Badge>
    </div>
  ),
};

// All Sizes
export const AllSizes: Story = {
  render: () => (
    <div className="flex items-center gap-3">
      <Badge size="sm" variant="primary">
        Small
      </Badge>
      <Badge size="md" variant="primary">
        Medium
      </Badge>
      <Badge size="lg" variant="primary">
        Large
      </Badge>
    </div>
  ),
};

// Language Labels (HaiLanGo specific)
export const LanguageLabels: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <Badge variant="primary">Russian</Badge>
      <Badge variant="secondary">Japanese</Badge>
      <Badge variant="info">Arabic</Badge>
      <Badge variant="warning">Persian</Badge>
      <Badge variant="danger">Hebrew</Badge>
    </div>
  ),
};

// Learning Status Badges
export const LearningStatus: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <Badge variant="success" dot>
        Completed
      </Badge>
      <Badge variant="primary" dot>
        In Progress
      </Badge>
      <Badge variant="warning" dot>
        Review Due
      </Badge>
      <Badge variant="danger" dot>
        Urgent Review
      </Badge>
    </div>
  ),
};

// Progress Badges
export const ProgressBadges: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <Badge variant="success">85%</Badge>
      <Badge variant="primary">45%</Badge>
      <Badge variant="warning">20%</Badge>
      <Badge variant="danger">5%</Badge>
    </div>
  ),
};

// Badge Group
export const Group: Story = {
  render: () => (
    <BadgeGroup>
      <Badge variant="primary">Vocabulary</Badge>
      <Badge variant="secondary">Grammar</Badge>
      <Badge variant="info">Listening</Badge>
      <Badge variant="success">Speaking</Badge>
    </BadgeGroup>
  ),
};

// Book Category Badges
export const BookCategories: Story = {
  render: () => (
    <div className="p-4 bg-white rounded-xl shadow-soft space-y-4">
      <h3 className="font-semibold text-text-primary">Russian Language 101</h3>
      <BadgeGroup>
        <Badge variant="primary" size="md">
          Beginner
        </Badge>
        <Badge variant="secondary" size="md">
          Vocabulary
        </Badge>
        <Badge variant="info" size="md">
          150 pages
        </Badge>
      </BadgeGroup>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Streak Badge
export const StreakBadge: Story = {
  render: () => (
    <Badge
      variant="warning"
      size="lg"
      icon={
        <svg aria-hidden="true" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zM7 9c0-2.76 2.24-5 5-5s5 2.24 5 5c0 2.88-2.88 7.19-5 9.88C9.92 16.21 7 11.85 7 9z" />
        </svg>
      }
    >
      7 Day Streak
    </Badge>
  ),
};

// Achievement Badges
export const AchievementBadges: Story = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft space-y-4">
      <h3 className="font-semibold text-text-primary">Achievements</h3>
      <div className="grid grid-cols-2 gap-3">
        <div className="flex items-center gap-2 p-3 bg-success/5 rounded-lg">
          <span className="text-2xl">100</span>
          <Badge variant="success" size="sm">
            Words Learned
          </Badge>
        </div>
        <div className="flex items-center gap-2 p-3 bg-primary/5 rounded-lg">
          <span className="text-2xl">85</span>
          <Badge variant="primary" size="sm">
            Avg. Score
          </Badge>
        </div>
        <div className="flex items-center gap-2 p-3 bg-warning/5 rounded-lg">
          <span className="text-2xl">7</span>
          <Badge variant="warning" size="sm">
            Day Streak
          </Badge>
        </div>
        <div className="flex items-center gap-2 p-3 bg-secondary/5 rounded-lg">
          <span className="text-2xl">3</span>
          <Badge variant="secondary" size="sm">
            Books
          </Badge>
        </div>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};
