#!/usr/bin/env python3

import numpy as np
import json
import sys
import os
from dataclasses import dataclass
from typing import Dict

@dataclass
class Context:
    language: str
    file_type: str
    diff_size: int
    time_of_day: int
    
    def to_vector(self, suggestion_type: str) -> np.ndarray:
        features = np.zeros(10)
        
        lang_map = {
            'go': 0,
            'python': 1, 
            'javascript': 2,
            'typescript': 2,
        }
        lang_idx = lang_map.get(self.language, 3)
        features[lang_idx] = 1.0
        
        type_map = {
            'error-handling': 4,
            'security': 5,
            'performance': 6,
            'reliability': 7,
        }
        type_idx = type_map.get(suggestion_type, 8)
        features[type_idx] = 1.0
        
        features[9] = min(self.diff_size / 100.0, 1.0)
        
        return features

class LinUCB:
    def __init__(self, n_features: int = 10, alpha: float = 0.5):
        self.alpha = alpha
        self.n_features = n_features
        self.A = np.identity(n_features)
        self.b = np.zeros(n_features)
        self.total_predictions = 0
        self.total_updates = 0
    
    def predict(self, context: Context, suggestion_type: str) -> float:
        x = context.to_vector(suggestion_type)
        
        try:
            A_inv = np.linalg.inv(self.A)
        except np.linalg.LinAlgError:
            A_inv = np.linalg.inv(self.A + 1e-6 * np.identity(self.n_features))
        
        theta = A_inv @ self.b
        expected_reward = theta.T @ x
        uncertainty = np.sqrt(x.T @ A_inv @ x)
        ucb = expected_reward + self.alpha * uncertainty
        
        self.total_predictions += 1
        
        return float(ucb)
    
    def update(self, context: Context, suggestion_type: str, reward: float):
        x = context.to_vector(suggestion_type)
        self.A += np.outer(x, x)
        self.b += reward * x
        self.total_updates += 1
    
    def get_stats(self) -> Dict:
        return {
            'total_predictions': self.total_predictions,
            'total_updates': self.total_updates,
            'exploration_param': self.alpha,
        }
    
    def save(self, path: str):
        np.savez(
            path,
            A=self.A,
            b=self.b,
            alpha=self.alpha,
            total_predictions=self.total_predictions,
            total_updates=self.total_updates,
        )
    
    def load(self, path: str):
        if not os.path.exists(path):
            raise FileNotFoundError(f"Model file not found: {path}")
        
        data = np.load(path)
        self.A = data['A']
        self.b = data['b']
        self.alpha = float(data['alpha'])
        self.total_predictions = int(data['total_predictions'])
        self.total_updates = int(data['total_updates'])

def main():
    if len(sys.argv) < 2:
        print(json.dumps({'error': 'Usage: python linucb.py <command>'}))
        sys.exit(1)
    
    command = sys.argv[1]
    
    model_dir = os.path.join(os.getcwd(), '.cortex')
    os.makedirs(model_dir, exist_ok=True)
    model_path = os.path.join(model_dir, 'bandit.npz')
    
    bandit = LinUCB()
    
    try:
        if command == 'init':
            if not os.path.exists(model_path):
                bandit.save(model_path)
                print(json.dumps({'status': 'ok', 'message': 'Model initialized', 'path': model_path}))
            else:
                print(json.dumps({'status': 'ok', 'message': 'Model already exists', 'path': model_path}))
            sys.exit(0)
        
        try:
            bandit.load(model_path)
        except FileNotFoundError:
            pass
        
        if command == 'predict':
            data = json.load(sys.stdin)
            
            context = Context(**data['context'])
            suggestions = data['suggestions']
            
            scores = {}
            for sugg in suggestions:
                score = bandit.predict(context, sugg['type'])
                scores[sugg['id']] = score
            
            print(json.dumps(scores))
        
        elif command == 'update':
            data = json.load(sys.stdin)
            
            context = Context(**data['context'])
            suggestion_type = data['suggestion_type']
            action = data['action']
            
            reward = 1.0 if action == 'accept' else 0.0
            
            bandit.update(context, suggestion_type, reward)
            bandit.save(model_path)
            
            print(json.dumps({'status': 'ok', 'reward': reward}))
        
        elif command == 'stats':
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
