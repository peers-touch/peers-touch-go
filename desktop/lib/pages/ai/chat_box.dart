import 'package:desktop/widgets/ai/agent_selector.dart';
import 'package:desktop/widgets/ai/chat_input.dart';
import 'package:desktop/widgets/ai/search_assistant.dart';
import 'package:flutter/material.dart';

class ChatBox extends StatelessWidget {
  const ChatBox({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12.0),
        border: Border.all(color: Colors.grey.shade200),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min, // Important for bottom nav bar
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              AgentSelector(),
              SearchAssistant(),
            ],
          ),
          const SizedBox(height: 8.0),
          const Divider(),
          const ChatInput(),
        ],
      ),
    );
  }
}