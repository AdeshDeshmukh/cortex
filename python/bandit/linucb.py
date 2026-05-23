#!/usr/bin/env python3
"""
LinUCB (Linear Upper Confidence Bound) Contextual Bandit
For learning user preferences in code review suggestions.
"""

import numpy as np
import json
import sys
import os
from dataclasses import dataclass, asdict
from typing import Dict, List, Tuple


@dataclass
class Context:
    """
    Review context features.
    
    These are the "signals" that help predict if user will accept a suggestion.
    """
    language: str        # "go", "python", "javascript", etc.
    file_type: str       # Same as language for now
    diff_size: int       # Number of lines changed
    time_of_day: int     # 0-23 (future: learn time preferences)
    
    def to_vector(self, suggestion_type: str) -> np.ndarray:
        """
        Convert context + suggestion type to feature vector.
        
        Why?
          Machine learning needs numbers, not strings.
          This converts {language: "go", type: "error-handling"}
          into [0, 1, 0, 0, 1, 0, 0.5]
        
        Feature vector structure (10 dimensions):
          [0-3]: One-hot encoded language (go/python/js/other)
          [4-8]: One-hot encoded suggestion type
          [9]:   Normalized diff size
        """
        features = np.zeros(10)
        
        # One-hot encode language
        # Why one-hot?
        #   Can't do: go=1, python=2, javascript=3
        #   That implies: go < python < javascript (false!)
        #   One-hot: go=[1,0,0], python=[0,1,0], js=[0,0,1]
        lang_map = {
            'go': 0,
            'python': 1, 
            'javascript': 2,
            'typescript': 2,  # Treat TS same as JS
        }
        lang_idx = lang_map.get(self.language, 3)  # 3 = "other"
        features[lang_idx] = 1.0
        
        # One-hot encode suggestion type
        type_map = {
            'error-handling': 4,
            'security': 5,
            'performance': 6,
            'reliability': 7,
        }
        type_idx = type_map.get(suggestion_type, 8)  # 8 = "other"
        features[type_idx] = 1.0
        
        # Normalize diff size to [0, 1]
        # Why normalize?
        #   diff_size can be 5 or 500
        #   Without normalization: 500 dominates other features
        #   With: both become ~0.05 and ~1.0 (comparable scale)
        features[9] = min(self.diff_size / 100.0, 1.0)
        
        return features


class LinUCB:
    """
    Linear Upper Confidence Bound contextual bandit.
    
    Paper: "A Contextual-Bandit Approach to Personalized News Article Recommendation"
    Link: https://arxiv.org/abs/1003.0146
    
    The Math (simplified):
      We maintain: A matrix and b vector
      A = design matrix (feature correlations)
      b = response vector (weighted rewards)
      
      Prediction: θ = A^(-1) * b
      Confidence: UCB = θ^T * x + α * sqrt(x^T * A^(-1) * x)
      
      Higher UCB = More likely user will accept
    """
    
    def __init__(self, n_features: int = 10, alpha: float = 0.5):
        """
        Initialize LinUCB bandit.
        
        Args:
            n_features: Dimension of feature vector (10 in our case)
            alpha: Exploration parameter (higher = more exploration)
        
        Why alpha = 0.5?
          Too low (0.1): Only exploits known good types
          Too high (2.0): Keeps trying bad types
          0.5: Balanced exploration/exploitation
        """
        self.alpha = alpha
        self.n_features = n_features
        
        # A = identity matrix (start with no knowledge)
        # Why identity?
        #   Represents "I'm equally uncertain about all features"
        #   As we learn, A fills with correlations
        self.A = np.identity(n_features)
        
        # b = zero vector (no rewards observed yet)
        self.b = np.zeros(n_features)
        
        # Track statistics
        self.total_predictions = 0
        self.total_updates = 0
    
    def predict(self, context: Context, suggestion_type: str) -> float:
        """
        Predict acceptance probability with upper confidence bound.
        
        Returns:
            UCB score (0.0 to ~1.5)
            Higher = More likely to be accepted
        
        How it works:
          1. Convert context to feature vector x
          2. Compute θ (our best guess of user preferences)
          3. Add confidence bonus (explore uncertain options)
          4. Return UCB score
        """
        x = context.to_vector(suggestion_type)
        
        # Invert A matrix
        # Why invert?
        #   We need A^(-1) to solve: θ = A^(-1) * b
        #   This is solving a system of linear equations
        try:
            A_inv = np.linalg.inv(self.A)
        except np.linalg.LinAlgError:
            # Matrix is singular (shouldn't happen, but be safe)
            # Add small noise to diagonal (regularization)
            A_inv = np.linalg.inv(self.A + 1e-6 * np.identity(self.n_features))
        
        # Compute θ (estimated parameters)
        # θ tells us: "How much does each feature correlate with acceptance?"
        theta = A_inv @ self.b
        
        # Expected reward (exploitation term)
        expected_reward = theta.T @ x
        
        # Confidence bonus (exploration term)
        # sqrt(x^T * A^(-1) * x) = uncertainty about this specific context
        # If we've never seen this context → High uncertainty → Explore it
        # If we've seen it often → Low uncertainty → Trust our estimate
        uncertainty = np.sqrt(x.T @ A_inv @ x)
        
        # Upper confidence bound
        ucb = expected_reward + self.alpha * uncertainty
        
        self.total_predictions += 1
        
        return float(ucb)
    
    def update(self, context: Context, suggestion_type: str, reward: float):
        """
        Update model with feedback.
        
        Args:
            context: Review context
            suggestion_type: Type of suggestion shown
            reward: 1.0 if accepted, 0.0 if rejected
        
        How it works:
          1. Get feature vector x for this context
          2. Update A: A = A + x*x^T (track feature correlations)
          3. Update b: b = b + reward*x (track weighted rewards)
          
        Why this formula?
          It's solving ridge regression incrementally.
          Each update refines our estimate of user preferences.
        """
        x = context.to_vector(suggestion_type)
        
        # Update A (design matrix)
        # A accumulates: "How often do we see each feature combination?"
        self.A += np.outer(x, x)
        
        # Update b (response vector)
        # b accumulates: "What rewards did we get for each feature?"
        self.b += reward * x
        
        self.total_updates += 1
    
    def get_stats(self) -> Dict:
        """Get bandit statistics."""
        return {
            'total_predictions': self.total_predictions,
            'total_updates': self.total_updates,
            'exploration_param': self.alpha,
        }
    
    def save(self, path: str):
        """
        Persist model to disk.
        
        Saves A and b matrices so learning survives program restarts.
        """
        np.savez(
            path,
            A=self.A,
            b=self.b,
            alpha=self.alpha,
            total_predictions=self.total_predictions,
            total_updates=self.total_updates,
        )
    
    def load(self, path: str):
        """Load model from disk."""
        if not os.path.exists(path):
            raise FileNotFoundError(f"Model file not found: {path}")
        
        data = np.load(path)
        self.A = data['A']
        self.b = data['b']
        self.alpha = float(data['alpha'])
        self.total_predictions = int(data['total_predictions'])
        self.total_updates = int(data['total_updates'])


def main():
    """
    CLI interface for Go integration.
    
    Commands:
      predict: Rank suggestions by predicted acceptance
      update:  Update model with user feedback
      stats:   Get model statistics
    
    Communication:
      Input: JSON via stdin
      Output: JSON via stdout
      Errors: stderr
    """
    if len(sys.argv) < 2:
        print(json.dumps({'error': 'Usage: python linucb.py <command>'}))
        sys.exit(1)
    
    command = sys.argv[1]
    
    # Model persistence path
    model_dir = os.path.join(os.getcwd(), '.cortex')
    os.makedirs(model_dir, exist_ok=True)
    model_path = os.path.join(model_dir, 'bandit.npz')
    
    # Initialize or load bandit
    bandit = LinUCB()
    try:
        bandit.load(model_path)
    except FileNotFoundError:
        # First run, no saved model yet
        pass
    
    try:
        if command == 'predict':
            # Read input from stdin
            data = json.load(sys.stdin)
            
            context = Context(**data['context'])
            suggestions = data['suggestions']
            
            # Score each suggestion
            scores = {}
            for sugg in suggestions:
                score = bandit.predict(context, sugg['type'])
                scores[sugg['id']] = score
            
            # Return scores
            print(json.dumps(scores))
        
        elif command == 'update':
            # Read feedback from stdin
            data = json.load(sys.stdin)
            
            context = Context(**data['context'])
            suggestion_type = data['suggestion_type']
            action = data['action']
            
            # Convert action to reward
            reward = 1.0 if action == 'accept' else 0.0
            
            # Update model
            bandit.update(context, suggestion_type, reward)
            
            # Save updated model
            bandit.save(model_path)
            
            print(json.dumps({'status': 'ok', 'reward': reward}))
        
        elif command == 'stats':
            # Return model statistics
            stats = bandit.get_stats()
            print(json.dumps(stats))
        
        else:
            print(json.dumps({'error': f'Unknown command: {command}'}))
            sys.exit(1)
    
    except Exception as e:
        print(json.dumps({'error': str(e)}), file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    main()