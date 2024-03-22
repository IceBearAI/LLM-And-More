import {
  IconAlertCircle,
  IconApps,
  IconBorderAll,
  IconBrandTabler,
  IconCircleDot,
  IconClipboard,
  IconFileDescription,
  IconHome,
  IconLogin,
  IconRotate,
  IconSettings,
  IconUserPlus,
  IconZoomCode,
  IconColumns,
  IconRowInsertBottom,
  IconEyeTable,
  IconSortAscending,
  IconPageBreak,
  IconFilter,
  IconBoxModel,
  IconServer,
  IconBrandAirtable,
  IconVolume2
} from "@tabler/icons-vue";

export interface menu {
  header?: string;
  title?: string;
  icon?: any;
  to?: string;
  divider?: boolean;
  chip?: string;
  chipColor?: string;
  chipVariant?: string;
  chipIcon?: string;
  children?: menu[];
  disabled?: boolean;
  subCaption?: string;
  class?: string;
  extraclass?: string;
  type?: string;
}

const horizontalItems: menu[] = [
  {
    title: "智能声纹",
    icon: IconHome,
    to: "#",
    children: [
      {
        title: "声纹识别",
        icon: IconVolume2,
        to: "/voice-print/analysis"
      }
    ]
  },
  {
    title: "Dashboard",
    icon: IconHome,
    to: "#",
    children: [
      {
        title: "Analytical",
        icon: IconCircleDot,
        to: "/dashboards/analytical"
      },
      {
        title: "eCommerce",
        icon: IconCircleDot,
        to: "/dashboards/ecommerce"
      },
      {
        title: "Modern",
        icon: IconCircleDot,
        to: "/dashboards/modern"
      }
    ]
  },
  {
    title: "Apps",
    icon: IconApps,
    to: "#",
    children: [
      {
        title: "Chats",
        icon: IconCircleDot,
        to: "/apps/chats"
      },
      {
        title: "Blog",
        icon: IconCircleDot,
        to: "/",
        children: [
          {
            title: "Posts",
            icon: IconCircleDot,
            to: "/apps/blog/posts"
          },
          {
            title: "Detail",
            icon: IconCircleDot,
            to: "/apps/blog/early-black-friday-amazon-deals-cheap-tvs-headphones"
          }
        ]
      },
      {
        title: "E-Commerce",
        icon: IconCircleDot,
        to: "/ecommerce/",
        children: [
          {
            title: "Shop",
            icon: IconCircleDot,
            to: "/ecommerce/products"
          },
          {
            title: "Detail",
            icon: IconCircleDot,
            to: "/ecommerce/product/detail/1"
          },
          {
            title: "List",
            icon: IconCircleDot,
            to: "/ecommerce/productlist"
          },
          {
            title: "Checkout",
            icon: IconCircleDot,
            to: "/ecommerce/checkout"
          }
        ]
      },
      {
        title: "User Profile",
        icon: IconCircleDot,
        to: "/",
        children: [
          {
            title: "Profile",
            icon: IconCircleDot,
            to: "/apps/user/profile"
          },
          {
            title: "Followers",
            icon: IconCircleDot,
            to: "/apps/user/profile/followers"
          },
          {
            title: "Friends",
            icon: IconCircleDot,
            to: "/apps/user/profile/friends"
          },
          {
            title: "Gallery",
            icon: IconCircleDot,
            to: "/apps/user/profile/gallery"
          }
        ]
      },
      {
        title: "Notes",
        icon: IconCircleDot,
        to: "/apps/notes"
      },
      {
        title: "Calendar",
        icon: IconCircleDot,
        to: "/apps/calendar"
      },
      {
        title: "Kanban",
        icon: IconCircleDot,
        to: "/apps/kanban"
      }
    ]
  },

  {
    title: "Pages",
    icon: IconClipboard,
    to: "#",
    children: [
      {
        title: "Widget",
        icon: IconCircleDot,
        to: "/widget-card",
        children: [
          {
            title: "Cards",
            icon: IconCircleDot,
            to: "/widgets/cards"
          },
          {
            title: "Banners",
            icon: IconCircleDot,
            to: "/widgets/banners"
          },
          {
            title: "Charts",
            icon: IconCircleDot,
            to: "/widgets/charts"
          }
        ]
      },
      {
        title: "UI",
        icon: IconCircleDot,
        to: "#",
        children: [
          {
            title: "Alert",
            icon: IconCircleDot,
            to: "/ui-components/alert"
          },
          {
            title: "Accordion",
            icon: IconCircleDot,
            to: "/ui-components/accordion"
          },
          {
            title: "Avatar",
            icon: IconCircleDot,
            to: "/ui-components/avatar"
          },
          {
            title: "Chip",
            icon: IconCircleDot,
            to: "/ui-components/chip"
          },
          {
            title: "Dialog",
            icon: IconCircleDot,
            to: "/ui-components/dialogs"
          },
          {
            title: "List",
            icon: IconCircleDot,
            to: "/ui-components/list"
          },
          {
            title: "Menus",
            icon: IconCircleDot,
            to: "/ui-components/menus"
          },
          {
            title: "Rating",
            icon: IconCircleDot,
            to: "/ui-components/rating"
          },
          {
            title: "Tabs",
            icon: IconCircleDot,
            to: "/ui-components/tabs"
          },
          {
            title: "Tooltip",
            icon: IconCircleDot,
            to: "/ui-components/tooltip"
          },
          {
            title: "Typography",
            icon: IconCircleDot,
            to: "/ui-components/typography"
          }
        ]
      },
      {
        title: "Charts",
        icon: IconCircleDot,
        to: "#",
        children: [
          {
            title: "Line",
            icon: IconCircleDot,
            to: "/charts/line-chart"
          },
          {
            title: "Gredient",
            icon: IconCircleDot,
            to: "/charts/gredient-chart"
          },
          {
            title: "Area",
            icon: IconCircleDot,
            to: "/charts/area-chart"
          },
          {
            title: "Candlestick",
            icon: IconCircleDot,
            to: "/charts/candlestick-chart"
          },
          {
            title: "Column",
            icon: IconCircleDot,
            to: "/charts/column-chart"
          },
          {
            title: "Doughnut & Pie",
            icon: IconCircleDot,
            to: "/charts/doughnut-pie-chart"
          },
          {
            title: "Radialbar & Radar",
            icon: IconCircleDot,
            to: "/charts/radialbar-chart"
          }
        ]
      },
      {
        title: "Auth",
        icon: IconCircleDot,
        to: "#",
        children: [
          {
            title: "Error",
            icon: IconAlertCircle,
            to: "/auth/404"
          },
          {
            title: "Maintenance",
            icon: IconSettings,
            to: "/auth/maintenance"
          },
          {
            title: "Login",
            icon: IconLogin,
            to: "#",
            children: [
              {
                title: "Side Login",
                icon: IconCircleDot,
                to: "/auth/login"
              },
              {
                title: "Boxed Login",
                icon: IconCircleDot,
                to: "/auth/login2"
              }
            ]
          },
          {
            title: "Register",
            icon: IconUserPlus,
            to: "#",
            children: [
              {
                title: "Side Register",
                icon: IconCircleDot,
                to: "/auth/register"
              },
              {
                title: "Boxed Register",
                icon: IconCircleDot,
                to: "/auth/register2"
              }
            ]
          },
          {
            title: "Forgot Password",
            icon: IconRotate,
            to: "#",
            children: [
              {
                title: "Side Forgot Password",
                icon: IconCircleDot,
                to: "/auth/forgot-password"
              },
              {
                title: "Boxed Forgot Password",
                icon: IconCircleDot,
                to: "/auth/forgot-password2"
              }
            ]
          },
          {
            title: "Two Steps",
            icon: IconZoomCode,
            to: "#",
            children: [
              {
                title: "Side Two Steps",
                icon: IconCircleDot,
                to: "/auth/two-step"
              },
              {
                title: "Boxed Two Steps",
                icon: IconCircleDot,
                to: "/auth/two-step2"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    title: "Forms",
    icon: IconFileDescription,
    to: "#",
    children: [
      {
        title: "Form Elements",
        icon: IconCircleDot,
        to: "/components/",
        children: [
          {
            title: "Autocomplete",
            icon: IconCircleDot,
            to: "/forms/form-elements/autocomplete"
          },
          {
            title: "Combobox",
            icon: IconCircleDot,
            to: "/forms/form-elements/combobox"
          },
          {
            title: "Button",
            icon: IconCircleDot,
            to: "/forms/form-elements/button"
          },
          {
            title: "Checkbox",
            icon: IconCircleDot,
            to: "/forms/form-elements/checkbox"
          },
          {
            title: "Custom Inputs",
            icon: IconCircleDot,
            to: "/forms/form-elements/custominputs"
          },
          {
            title: "File Inputs",
            icon: IconCircleDot,
            to: "/forms/form-elements/fileinputs"
          },
          {
            title: "Radio",
            icon: IconCircleDot,
            to: "/forms/form-elements/radio"
          },
          {
            title: "Select",
            icon: IconCircleDot,
            to: "/forms/form-elements/select"
          },
          {
            title: "Date Time",
            icon: IconCircleDot,
            to: "/forms/form-elements/date-time"
          },
          {
            title: "Slider",
            icon: IconCircleDot,
            to: "/forms/form-elements/slider"
          },
          {
            title: "Switch",
            icon: IconCircleDot,
            to: "/forms/form-elements/switch"
          }
        ]
      },
      {
        title: "Form Layout",
        icon: IconCircleDot,
        to: "/forms/form-layouts"
      },
      {
        title: "Form Horizontal",
        icon: IconCircleDot,
        to: "/forms/form-horizontal"
      },
      {
        title: "Form Vertical",
        icon: IconCircleDot,
        to: "/forms/form-vertical"
      },
      {
        title: "Form Custom",
        icon: IconCircleDot,
        to: "/forms/form-custom"
      },
      {
        title: "Form Validation",
        icon: IconCircleDot,
        to: "/forms/form-validation"
      }
    ]
  },
  {
    title: "Tables",
    icon: IconBorderAll,
    to: "#",
    children: [
      {
        title: "Basic Table",
        icon: IconCircleDot,
        to: "/tables/basic"
      },
      {
        title: "Dark Table",
        icon: IconCircleDot,
        to: "/tables/dark"
      },
      {
        title: "Density Table",
        icon: IconCircleDot,
        to: "/tables/density"
      },
      {
        title: "Fixed Header Table",
        icon: IconCircleDot,
        to: "/tables/fixed-header"
      },
      {
        title: "Height Table",
        icon: IconCircleDot,
        to: "/tables/height"
      },
      {
        title: "Editable Table",
        icon: IconCircleDot,
        to: "/tables/editable"
      }
    ]
  },
  {
    title: "Data Tables",
    icon: IconBrandAirtable,
    to: "#",
    children: [
      {
        title: "Basic Table",
        icon: IconColumns,
        to: "/datatables/basic"
      },
      {
        title: "Header Table",
        icon: IconRowInsertBottom,
        to: "/datatables/header"
      },
      {
        title: "Selection Table",
        icon: IconEyeTable,
        to: "/datatables/selection"
      },
      {
        title: "Sorting Table",
        icon: IconSortAscending,
        to: "/datatables/sorting"
      },
      {
        title: "Pagination Table",
        icon: IconPageBreak,
        to: "/datatables/pagination"
      },
      {
        title: "Filtering Table",
        icon: IconFilter,
        to: "/datatables/filtering"
      },
      {
        title: "Grouping Table",
        icon: IconBoxModel,
        to: "/datatables/grouping"
      },
      {
        title: "Table Slots",
        icon: IconServer,
        to: "/datatables/slots"
      }
    ]
  },
  {
    title: "Icons",
    icon: IconBrandTabler,
    to: "#",
    children: [
      {
        title: "Material",
        icon: IconCircleDot,
        to: "/icons/material"
      },
      {
        title: "Tabler",
        icon: IconCircleDot,
        to: "/icons/tabler"
      }
    ]
  }
];

export default horizontalItems;
