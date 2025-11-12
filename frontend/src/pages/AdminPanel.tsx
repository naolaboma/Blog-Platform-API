import React, { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { adminAPI } from '@/lib/auth';
import { User } from '@/types';
import {
  Users,
  Shield,
  Crown,
  User as UserIcon,
  Mail,
  Calendar,
  ChevronUp,
  ChevronDown,
} from 'lucide-react';
import toast from 'react-hot-toast';

interface UserWithRole extends User {
  isLoading?: boolean;
}

const AdminPanel: React.FC = () => {
  const { user } = useAuth();
  const [users, setUsers] = useState<UserWithRole[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [sortField, setSortField] = useState<'createdAt' | 'username'>('createdAt');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc');

  useEffect(() => {
    if (!user || user.role !== 'admin') {
      return;
    }
    fetchUsers();
  }, [user]);

  const fetchUsers = async () => {
    setIsLoading(true);
    try {
      const data = await adminAPI.getAllUsers();
      setUsers(data.users || []);
    } catch (error) {
      console.error('Error fetching users:', error);
      toast.error('Failed to load users');
    } finally {
      setIsLoading(false);
    }
  };

  const handleRoleChange = async (userId: string, newRole: 'admin' | 'user') => {
    setUsers(prevUsers =>
      prevUsers.map(u =>
        u.id === userId ? { ...u, isLoading: true } : u
      )
    );

    try {
      const currentUser = users.find(u => u.id === userId);
      if (!currentUser) return;

      if (currentUser.role === 'admin' && newRole === 'user') {
        await adminAPI.demoteUser(userId, newRole);
        toast.success('User demoted successfully');
      } else if (currentUser.role === 'user' && newRole === 'admin') {
        await adminAPI.promoteUser(userId, newRole);
        toast.success('User promoted successfully');
      }

      setUsers(prevUsers =>
        prevUsers.map(u =>
          u.id === userId ? { ...u, role: newRole, isLoading: false } : u
        )
      );
    } catch (error) {
      console.error('Error changing user role:', error);
      toast.error('Failed to change user role');
      setUsers(prevUsers =>
        prevUsers.map(u =>
          u.id === userId ? { ...u, isLoading: false } : u
        )
      );
    }
  };

  const handleSort = (field: 'createdAt' | 'username') => {
    if (sortField === field) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('asc');
    }
  };

  const sortedUsers = [...users].sort((a, b) => {
    let comparison = 0;
    
    if (sortField === 'createdAt') {
      comparison = new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
    } else if (sortField === 'username') {
      comparison = a.username.localeCompare(b.username);
    }
    
    return sortDirection === 'asc' ? comparison : -comparison;
  });

  if (!user || user.role !== 'admin') {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Shield className="mx-auto h-12 w-12 text-gray-400" />
          <h2 className="mt-2 text-lg font-medium text-gray-900">Access Denied</h2>
          <p className="mt-1 text-sm text-gray-500">
            You don't have permission to access this page.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 flex items-center">
            <Crown className="mr-3 text-yellow-500" />
            Admin Panel
          </h1>
          <p className="mt-2 text-gray-600">
            Manage user roles and permissions
          </p>
        </div>

        <div className="card">
          <div className="px-6 py-4 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-medium text-gray-900 flex items-center">
                <Users className="mr-2 text-gray-400" />
                User Management
              </h2>
              <div className="text-sm text-gray-500">
                Total Users: {users.length}
              </div>
            </div>
          </div>

          {isLoading ? (
            <div className="flex justify-center items-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      User
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                        onClick={() => handleSort('username')}>
                      <div className="flex items-center space-x-1">
                        <span>Username</span>
                        {sortField === 'username' && (
                          sortDirection === 'asc' ? <ChevronUp size={14} /> : <ChevronDown size={14} />
                        )}
                      </div>
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Email
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                        onClick={() => handleSort('createdAt')}>
                      <div className="flex items-center space-x-1">
                        <span>Joined</span>
                        {sortField === 'createdAt' && (
                          sortDirection === 'asc' ? <ChevronUp size={14} /> : <ChevronDown size={14} />
                        )}
                      </div>
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Role
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {sortedUsers.map((userItem) => (
                    <tr key={userItem.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center">
                          {userItem.profilePicture ? (
                            <img
                              className="h-10 w-10 rounded-full object-cover"
                              src={userItem.profilePicture.filePath}
                              alt={userItem.username}
                            />
                          ) : (
                            <div className="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center">
                              <UserIcon size={20} className="text-primary-600" />
                            </div>
                          )}
                          <div className="ml-4">
                            <div className="text-sm font-medium text-gray-900">
                              {userItem.username}
                            </div>
                            <div className="text-sm text-gray-500">
                              {userItem.emailVerified ? 'Verified' : 'Not Verified'}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {userItem.username}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        <div className="flex items-center">
                          <Mail size={14} className="mr-1" />
                          {userItem.email}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        <div className="flex items-center">
                          <Calendar size={14} className="mr-1" />
                          {new Date(userItem.createdAt).toLocaleDateString()}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                          userItem.role === 'admin'
                            ? 'bg-purple-100 text-purple-800'
                            : 'bg-green-100 text-green-800'
                        }`}>
                          {userItem.role === 'admin' ? (
                            <><Crown size={12} className="mr-1" /> Admin</>
                          ) : (
                            <><UserIcon size={12} className="mr-1" /> User</>
                          )}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        {userItem.id !== user.id && (
                          <div className="flex space-x-2">
                            {userItem.role === 'user' ? (
                              <button
                                onClick={() => handleRoleChange(userItem.id, 'admin')}
                                disabled={userItem.isLoading}
                                className="btn btn-primary text-xs disabled:opacity-50"
                              >
                                Promote
                              </button>
                            ) : (
                              <button
                                onClick={() => handleRoleChange(userItem.id, 'user')}
                                disabled={userItem.isLoading}
                                className="btn btn-secondary text-xs disabled:opacity-50"
                              >
                                Demote
                              </button>
                            )}
                          </div>
                        )}
                        {userItem.id === user.id && (
                          <span className="text-xs text-gray-500">Current User</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>

              {users.length === 0 && (
                <div className="text-center py-12">
                  <Users className="mx-auto h-12 w-12 text-gray-400" />
                  <h3 className="mt-2 text-sm font-medium text-gray-900">No users found</h3>
                  <p className="mt-1 text-sm text-gray-500">
                    No users have registered yet.
                  </p>
                </div>
              )}
            </div>
          )}
        </div>

        <div className="mt-8 bg-yellow-50 border border-yellow-200 rounded-lg p-6">
          <div className="flex">
            <div className="flex-shrink-0">
              <Shield className="h-6 w-6 text-yellow-600" />
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-yellow-800">Admin Guidelines</h3>
              <div className="mt-2 text-sm text-yellow-700">
                <ul className="list-disc list-inside space-y-1">
                  <li>Promote trusted users to admin role for additional permissions</li>
                  <li>Admins can manage all blogs and user accounts</li>
                  <li>Demote admins if they no longer need elevated access</li>
                  <li>You cannot change your own role - contact another admin if needed</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminPanel;
