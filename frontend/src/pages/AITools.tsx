import React, { useState } from 'react';
import { aiAPI } from '@/lib/ai';
import { Sparkles, FileText, Wand2, Lightbulb } from 'lucide-react';
import toast from 'react-hot-toast';

const AITools: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'generate' | 'enhance' | 'ideas'>('generate');
  const [topic, setTopic] = useState('');
  const [content, setContent] = useState('');
  const [keywords, setKeywords] = useState('');
  const [generatedContent, setGeneratedContent] = useState('');
  const [enhancedContent, setEnhancedContent] = useState('');
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [isGenerating, setIsGenerating] = useState(false);
  const [isEnhancing, setIsEnhancing] = useState(false);
  const [isGettingIdeas, setIsGettingIdeas] = useState(false);

  const handleGenerate = async () => {
    if (!topic.trim()) {
      toast.error('Please enter a topic');
      return;
    }

    setIsGenerating(true);
    try {
      const response = await aiAPI.generateBlog({ topic: topic.trim() });
      setGeneratedContent(response.content);
      toast.success('Content generated successfully!');
    } catch (error) {
      console.error('Error generating content:', error);
      toast.error('Failed to generate content');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleEnhance = async () => {
    if (!content.trim()) {
      toast.error('Please enter some content to enhance');
      return;
    }

    setIsEnhancing(true);
    try {
      const response = await aiAPI.enhanceBlog({ content: content.trim() });
      setEnhancedContent(response.content);
      toast.success('Content enhanced successfully!');
    } catch (error) {
      console.error('Error enhancing content:', error);
      toast.error('Failed to enhance content');
    } finally {
      setIsEnhancing(false);
    }
  };

  const handleGetIdeas = async () => {
    const keywordArray = keywords
      .split(',')
      .map(k => k.trim())
      .filter(k => k.length > 0);

    if (keywordArray.length === 0) {
      toast.error('Please enter some keywords');
      return;
    }

    setIsGettingIdeas(true);
    try {
      const response = await aiAPI.suggestIdeas({ keywords: keywordArray });
      setSuggestions(response.ideas);
      toast.success('Ideas generated successfully!');
    } catch (error) {
      console.error('Error getting ideas:', error);
      toast.error('Failed to get ideas');
    } finally {
      setIsGettingIdeas(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toast.success('Copied to clipboard!');
  };

  const tabs = [
    { id: 'generate', label: 'Generate Blog', icon: FileText },
    { id: 'enhance', label: 'Enhance Content', icon: Wand2 },
    { id: 'ideas', label: 'Get Ideas', icon: Lightbulb },
  ];

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">
            AI Writing Tools
          </h1>
          <p className="text-xl text-gray-600">
            Leverage AI to create, enhance, and get inspired for your blog content
          </p>
        </div>

        <div className="card">
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8 px-6" aria-label="Tabs">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id as any)}
                  className={`py-4 px-1 border-b-2 font-medium text-sm flex items-center space-x-2 ${
                    activeTab === tab.id
                      ? 'border-primary-500 text-primary-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  <tab.icon size={16} />
                  <span>{tab.label}</span>
                </button>
              ))}
            </nav>
          </div>

          <div className="p-6">
            {activeTab === 'generate' && (
              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Topic
                  </label>
                  <input
                    type="text"
                    value={topic}
                    onChange={(e) => setTopic(e.target.value)}
                    className="input"
                    placeholder="Enter a topic for your blog (e.g., 'The benefits of remote work')"
                  />
                </div>

                <button
                  onClick={handleGenerate}
                  disabled={isGenerating}
                  className="btn btn-primary disabled:opacity-50"
                >
                  <Sparkles size={16} className="mr-2" />
                  {isGenerating ? 'Generating...' : 'Generate Content'}
                </button>

                {generatedContent && (
                  <div className="mt-6">
                    <div className="flex justify-between items-center mb-2">
                      <h3 className="text-lg font-medium text-gray-900">Generated Content</h3>
                      <button
                        onClick={() => copyToClipboard(generatedContent)}
                        className="btn btn-secondary text-sm"
                      >
                        Copy to Clipboard
                      </button>
                    </div>
                    <div className="bg-gray-50 p-4 rounded-lg">
                      <pre className="whitespace-pre-wrap text-sm text-gray-700">{generatedContent}</pre>
                    </div>
                  </div>
                )}
              </div>
            )}

            {activeTab === 'enhance' && (
              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Original Content
                  </label>
                  <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    rows={8}
                    className="input resize-none"
                    placeholder="Paste your content here to enhance it with AI..."
                  />
                </div>

                <button
                  onClick={handleEnhance}
                  disabled={isEnhancing}
                  className="btn btn-primary disabled:opacity-50"
                >
                  <Wand2 size={16} className="mr-2" />
                  {isEnhancing ? 'Enhancing...' : 'Enhance Content'}
                </button>

                {enhancedContent && (
                  <div className="mt-6">
                    <div className="flex justify-between items-center mb-2">
                      <h3 className="text-lg font-medium text-gray-900">Enhanced Content</h3>
                      <button
                        onClick={() => copyToClipboard(enhancedContent)}
                        className="btn btn-secondary text-sm"
                      >
                        Copy to Clipboard
                      </button>
                    </div>
                    <div className="bg-gray-50 p-4 rounded-lg">
                      <pre className="whitespace-pre-wrap text-sm text-gray-700">{enhancedContent}</pre>
                    </div>
                  </div>
                )}
              </div>
            )}

            {activeTab === 'ideas' && (
              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Keywords
                  </label>
                  <input
                    type="text"
                    value={keywords}
                    onChange={(e) => setKeywords(e.target.value)}
                    className="input"
                    placeholder="Enter keywords separated by commas (e.g., technology, programming, AI)"
                  />
                  <p className="mt-1 text-sm text-gray-500">
                    Provide keywords to get blog topic suggestions
                  </p>
                </div>

                <button
                  onClick={handleGetIdeas}
                  disabled={isGettingIdeas}
                  className="btn btn-primary disabled:opacity-50"
                >
                  <Lightbulb size={16} className="mr-2" />
                  {isGettingIdeas ? 'Getting Ideas...' : 'Get Ideas'}
                </button>

                {suggestions.length > 0 && (
                  <div className="mt-6">
                    <h3 className="text-lg font-medium text-gray-900 mb-4">Blog Ideas</h3>
                    <div className="space-y-3">
                      {suggestions.map((idea, index) => (
                        <div key={index} className="bg-gray-50 p-4 rounded-lg">
                          <div className="flex justify-between items-start">
                            <p className="text-sm text-gray-700 flex-1">{idea}</p>
                            <button
                              onClick={() => copyToClipboard(idea)}
                              className="ml-4 btn btn-secondary text-xs"
                            >
                              Copy
                            </button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        <div className="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
          <div className="flex">
            <div className="flex-shrink-0">
              <Sparkles className="h-6 w-6 text-blue-600" />
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">AI Tips</h3>
              <div className="mt-2 text-sm text-blue-700">
                <ul className="list-disc list-inside space-y-1">
                  <li>Be specific with your topics for better results</li>
                  <li>Use enhancement to improve readability and flow</li>
                  <li>Combine multiple keywords for diverse ideas</li>
                  <li>Always review and edit AI-generated content</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AITools;
