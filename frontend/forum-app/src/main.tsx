import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './react/App';
import './react/styles.css';

const container = document.getElementById('root');

if (container !== null) {
  createRoot(container).render(
    React.createElement(BrowserRouter, null, React.createElement(App)),
  );
}
