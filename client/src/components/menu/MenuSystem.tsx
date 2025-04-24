// src/components/menu/MenuSystem.tsx
import React, { useState, useEffect } from 'react';
import TreatModal from './TreatModal';
import InformationModal from './InformationModal';
import ActionPlanModal from './ActionPlanModal';
import ActivityLogModal from './ActivityLogModal';
import './MenuSystem.css';
import ESPModal from './ESPModal';

// Import the icons (adjust paths as needed)
var informationIconSrc = '/assets/icons/information-icon.png';
var treatsIconSrc = '/assets/icons/treats-icon.png';
var actionPlanIconSrc = '/assets/icons/action-plan-icon.png';
var activityLogIconSrc = '/assets/icons/activity-log-icon.png';

interface MenuSystemProps {
  selectedEntity: { type: string; [key: string]: any } | null;
  onClose: () => void;
}

interface MenuOption {
  id: string;
  icon: string; // Fallback text
  iconSrc: string; // Image URL
  component: React.ReactNode;
}

const MenuSystem: React.FC<MenuSystemProps> = ({ selectedEntity, onClose }) => {
  const [activeMenu, setActiveMenu] = useState<string | null>(null);

  // Automatically open the Information modal when a new entity is selected
  useEffect(() => {
    if (selectedEntity) {
      setActiveMenu('information');
    } else {
      setActiveMenu(null);
    }
  }, [selectedEntity]);

  if (!selectedEntity) return null;

  // Define menu options based on entity type
  const menuOptions: MenuOption[] = selectedEntity.type === 'gotchi'
    ? [
        { id: 'information', icon: '?', iconSrc: informationIconSrc, component: <InformationModal entity={selectedEntity} /> },
        { id: 'treats', icon: 'i', iconSrc: treatsIconSrc, component: <TreatModal entity={selectedEntity}/> },
        { id: 'action-plan', icon: 'A', iconSrc: actionPlanIconSrc, component: <ActionPlanModal entity={selectedEntity} /> },
        { id: 'activity-log', icon: 'L', iconSrc: activityLogIconSrc, component: <ActivityLogModal entity={selectedEntity} /> },
      ]
    : [
        { id: 'information', icon: '?', iconSrc: informationIconSrc, component: <InformationModal entity={selectedEntity} /> },
      ];

  const handleMenuClick = (menuId: string) => {
    setActiveMenu(menuId === activeMenu ? null : menuId);
  };

  return (
    <div className="menu-system">
      <div className="menu-icons">
        {menuOptions.map((option) => (
          <div
            key={option.id}
            className={`menu-icon ${activeMenu === option.id ? 'active' : ''}`}
            onClick={() => handleMenuClick(option.id)}
          >
            <img
              src={option.iconSrc}
              alt={option.id}
              className="menu-icon-image"
              onError={(e) => {
                const imgElement = e.currentTarget;
                imgElement.style.display = 'none'; // Hide the broken image
                const nextSibling = imgElement.nextSibling;
                // Ensure nextSibling exists and is an HTMLElement before accessing style
                if (nextSibling && nextSibling instanceof HTMLElement) {
                  nextSibling.style.display = 'flex'; // Show the fallback text
                }
              }}
            />
            <span
              className="menu-icon-fallback"
              style={{ display: 'none' }} // Hidden by default, shown on image error
            >
              {option.icon}
            </span>
          </div>
        ))}
      </div>

      {activeMenu && (
        <div className="menu-page">
          {menuOptions.find((option) => option.id === activeMenu)?.component}
        </div>
      )}
      {/* {selectedEntity.type === "gotchi" && 
      <ESPModal entity={selectedEntity}/>
      }  */}
    </div>
  );
};

export default MenuSystem;
