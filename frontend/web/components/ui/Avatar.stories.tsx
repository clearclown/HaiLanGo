import type { Meta, StoryObj } from '@storybook/react';
import { Avatar, AvatarGroup } from './Avatar';

const meta = {
  title: 'UI/Avatar',
  component: Avatar,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component: 'Avatar component for user profile images, initials, and status indicators.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    size: {
      control: 'select',
      options: ['xs', 'sm', 'md', 'lg', 'xl', '2xl'],
      description: 'Avatar size',
    },
    variant: {
      control: 'select',
      options: ['circle', 'rounded'],
      description: 'Avatar shape variant',
    },
    status: {
      control: 'select',
      options: ['online', 'offline', 'busy', 'away'],
      description: 'User status indicator',
    },
    ring: {
      control: 'boolean',
      description: 'Show ring around avatar',
    },
  },
} satisfies Meta<typeof Avatar>;

export default meta;
type Story = StoryObj<typeof meta>;

// Default Avatar
export const Default: Story = {
  args: {
    alt: 'User',
    fallback: 'U',
    size: 'md',
  },
};

// With Image
export const WithImage: Story = {
  args: {
    src: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop',
    alt: 'John Doe',
    size: 'md',
  },
};

// With Initials
export const WithInitials: Story = {
  args: {
    alt: 'Taro Yamada',
    size: 'md',
  },
};

// With Status
export const WithStatus: Story = {
  args: {
    src: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop',
    alt: 'User',
    status: 'online',
    size: 'md',
  },
};

// All Sizes
export const AllSizes: Story = {
  render: () => (
    <div className="flex items-end gap-4">
      <Avatar size="xs" alt="XS" />
      <Avatar size="sm" alt="SM" />
      <Avatar size="md" alt="MD" />
      <Avatar size="lg" alt="LG" />
      <Avatar size="xl" alt="XL" />
      <Avatar size="2xl" alt="2XL" />
    </div>
  ),
};

// All Status Types
export const AllStatusTypes: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <div className="text-center">
        <Avatar
          src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
          alt="User"
          status="online"
          size="lg"
        />
        <p className="text-xs text-text-secondary mt-2">Online</p>
      </div>
      <div className="text-center">
        <Avatar
          src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
          alt="User"
          status="away"
          size="lg"
        />
        <p className="text-xs text-text-secondary mt-2">Away</p>
      </div>
      <div className="text-center">
        <Avatar
          src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
          alt="User"
          status="busy"
          size="lg"
        />
        <p className="text-xs text-text-secondary mt-2">Busy</p>
      </div>
      <div className="text-center">
        <Avatar
          src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
          alt="User"
          status="offline"
          size="lg"
        />
        <p className="text-xs text-text-secondary mt-2">Offline</p>
      </div>
    </div>
  ),
};

// Rounded Variant
export const RoundedVariant: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <Avatar
        src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
        alt="User"
        variant="circle"
        size="xl"
      />
      <Avatar
        src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
        alt="User"
        variant="rounded"
        size="xl"
      />
    </div>
  ),
};

// With Ring
export const WithRing: Story = {
  args: {
    src: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop',
    alt: 'User',
    ring: true,
    size: 'xl',
  },
};

// Avatar Group
export const Group: Story = {
  render: () => (
    <AvatarGroup max={4} size="md">
      <Avatar
        src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
        alt="User 1"
      />
      <Avatar
        src="https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=100&h=100&fit=crop"
        alt="User 2"
      />
      <Avatar
        src="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop"
        alt="User 3"
      />
      <Avatar
        src="https://images.unsplash.com/photo-1527980965255-d3b416303d12?w=100&h=100&fit=crop"
        alt="User 4"
      />
      <Avatar alt="User 5" />
      <Avatar alt="User 6" />
    </AvatarGroup>
  ),
};

// User Profile Card (HaiLanGo specific)
export const UserProfileCard: Story = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft flex items-center gap-4">
      <Avatar
        src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
        alt="Taro Yamada"
        size="xl"
        status="online"
        ring
      />
      <div>
        <h3 className="font-semibold text-text-primary">山田 太郎</h3>
        <p className="text-sm text-text-secondary">Level 12 • Premium</p>
        <div className="flex items-center gap-2 mt-1">
          <span className="text-xs text-warning">🔥 7 day streak</span>
        </div>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Learning Leaderboard
export const LearningLeaderboard: Story = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft space-y-4 w-80">
      <h3 className="font-semibold text-text-primary">This Week's Top Learners</h3>
      <div className="space-y-3">
        <div className="flex items-center gap-3">
          <span className="w-6 text-center font-bold text-warning">1</span>
          <Avatar
            src="https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=100&h=100&fit=crop"
            alt="User 1"
            size="sm"
          />
          <div className="flex-1">
            <p className="text-sm font-medium text-text-primary">佐藤 花子</p>
            <p className="text-xs text-text-secondary">12.5 hours</p>
          </div>
          <span className="text-lg">🥇</span>
        </div>
        <div className="flex items-center gap-3">
          <span className="w-6 text-center font-bold text-text-secondary">2</span>
          <Avatar
            src="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop"
            alt="User 2"
            size="sm"
          />
          <div className="flex-1">
            <p className="text-sm font-medium text-text-primary">田中 一郎</p>
            <p className="text-xs text-text-secondary">10.2 hours</p>
          </div>
          <span className="text-lg">🥈</span>
        </div>
        <div className="flex items-center gap-3">
          <span className="w-6 text-center font-bold text-text-secondary">3</span>
          <Avatar
            src="https://images.unsplash.com/photo-1527980965255-d3b416303d12?w=100&h=100&fit=crop"
            alt="User 3"
            size="sm"
          />
          <div className="flex-1">
            <p className="text-sm font-medium text-text-primary">鈴木 次郎</p>
            <p className="text-xs text-text-secondary">8.7 hours</p>
          </div>
          <span className="text-lg">🥉</span>
        </div>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Study Group
export const StudyGroup: Story = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft w-80">
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-semibold text-text-primary">Russian Study Group</h3>
        <span className="text-xs text-text-secondary">5 members</span>
      </div>
      <AvatarGroup max={5} size="md">
        <Avatar
          src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
          alt="Member 1"
          status="online"
        />
        <Avatar
          src="https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=100&h=100&fit=crop"
          alt="Member 2"
          status="online"
        />
        <Avatar
          src="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop"
          alt="Member 3"
          status="away"
        />
        <Avatar alt="Member 4" status="offline" />
        <Avatar alt="Member 5" status="offline" />
      </AvatarGroup>
      <p className="text-sm text-text-secondary mt-3">2 members currently online</p>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Settings Profile Section
export const SettingsProfile: Story = {
  render: () => (
    <div className="p-6 bg-white rounded-xl shadow-soft w-96">
      <h3 className="font-semibold text-text-primary mb-6">Profile Settings</h3>
      <div className="flex items-center gap-6">
        <div className="relative">
          <Avatar
            src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop"
            alt="Your Profile"
            size="2xl"
          />
          <button
            type="button"
            aria-label="Change profile photo"
            className="absolute bottom-0 right-0 w-8 h-8 bg-primary text-white rounded-full flex items-center justify-center shadow-md hover:bg-primary/90 transition-colors"
          >
            <svg
              aria-hidden="true"
              className="w-4 h-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
          </button>
        </div>
        <div className="flex-1">
          <p className="text-sm text-text-secondary mb-1">Display Name</p>
          <p className="font-medium text-text-primary">山田 太郎</p>
          <p className="text-sm text-text-secondary mt-3 mb-1">Email</p>
          <p className="font-medium text-text-primary">taro@example.com</p>
        </div>
      </div>
    </div>
  ),
  parameters: {
    backgrounds: { default: 'gray' },
  },
};

// Fallback States
export const FallbackStates: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <div className="text-center">
        <Avatar alt="Taro Yamada" size="lg" />
        <p className="text-xs text-text-secondary mt-2">From name</p>
      </div>
      <div className="text-center">
        <Avatar fallback="TY" size="lg" />
        <p className="text-xs text-text-secondary mt-2">Custom initials</p>
      </div>
      <div className="text-center">
        <Avatar size="lg" />
        <p className="text-xs text-text-secondary mt-2">No data</p>
      </div>
      <div className="text-center">
        <Avatar src="/broken-image.jpg" alt="Broken" size="lg" />
        <p className="text-xs text-text-secondary mt-2">Broken image</p>
      </div>
    </div>
  ),
};
