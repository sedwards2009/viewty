TODO
====

* Style / color theme system
* Integrate Micro (Smidgen) editor widget
  * Input field widget
  * Input text area widget
  * Text editor widget
* Declaritive syntax / construction
* Scrollarea widget
* Menu bar
* List Selection widget
* Table view widget
* File dialog
* Form/dialog widget with key shortcuts and cursor navigation

Theme / Style use cases
=======================

"I want to make a special kind of input field. It is its own Widget but I want it to get the styles from the normal input field which I'm trying to extend."

"I want one central place where I can set the default background and foregroud text colours for all widgets."

"Should styles returned by the theme system depend on the state of the widget?"


Style JSON:

{
  "*": {
    "foregroundColor": "#ffffff",
    "backgroundColor": "#000000"
  },
  "InputField": {
    "foregroundColor": "#ffffff",
    "backgroundColor": "#000000"
  },
  "NumberInputField": {
    "from": "InputField"
  }
}

How should "classes" on Widget work?
